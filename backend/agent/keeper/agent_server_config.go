package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/key"
	"github.com/stephenzsy/small-kms/backend/managedapp"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentServerConfiguration interface {
	NextWaitInterval() time.Duration
	TLSCertificateBundleFile() string
	Version() string
	VerifyJWTKeys() []cloudkey.JsonWebKey
}

type agentServerConfiguration struct {
	TLSCertificateFile string           `json:"tlsCertificate"`
	JWTVerifyKeys      []key.JsonWebKey `json:"jwtVerifyKeys"`
	fetchedConfig      *managedapp.AgentConfigServer
}

// VerifyJWTKeys implements AgentServerConfiguration.
func (asc *agentServerConfiguration) VerifyJWTKeys() []cloudkey.JsonWebKey {
	return asc.JWTVerifyKeys
}

// TLSCertificateBundleFile implements AgentServerConfiguration.
func (asc *agentServerConfiguration) TLSCertificateBundleFile() string {
	return asc.TLSCertificateFile
}

func (asc *agentServerConfiguration) Version() string {
	return asc.fetchedConfig.Version
}

func (c *agentServerConfiguration) NextWaitInterval() time.Duration {
	d := time.Until(c.fetchedConfig.RefreshAfter)
	if d < time.Minute*5 {
		return time.Minute * 5
	}
	return d
}

type agentConfigServerProcessor struct {
	envConfig         *agentcommon.AgentEnv
	configDir         string
	readyConfig       *agentServerConfiguration
	configProvisioner configProvisioner
}

type configProvisioner struct {
	processor    *agentConfigServerProcessor
	versionedDir string
	config       *managedapp.AgentConfigServer
}

func (p *configProvisioner) provision(c context.Context) (*agentServerConfiguration, error) {
	logger := log.Ctx(c)
	if _, err := os.Stat(p.versionedDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(p.versionedDir, 0700)
		} else {
			return nil, err
		}
	}

	agentClient, err := p.processor.envConfig.AgentClient()
	if err != nil {
		return nil, err
	}

	if fileBytes, err := json.Marshal(p.config); err != nil {
		return nil, err
	} else if err := os.WriteFile(filepath.Join(p.versionedDir, "config.json"), fileBytes, 0600); err != nil {
		return nil, err
	}

	tlsCertFilePath := filepath.Join(p.versionedDir, "tls-cert.pem")
	if _, err := os.Stat(tlsCertFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// pull certificate
			resp, err := agentClient.GetCertificateWithResponse(c, base.NamespaceKindServicePrincipal, base.ID("me"), p.config.TlsCertificateId)
			if err != nil {
				return nil, err
			}
			cert := resp.JSON200
			azSecretsClient, err := p.processor.envConfig.AzSecretsClient()
			if err != nil {
				return nil, err
			}
			sid := azsecrets.ID(cert.KeyVaultSecretID)
			getSecretResposne, err := azSecretsClient.GetSecret(c, sid.Name(), sid.Version(), nil)
			if err != nil {
				return nil, err
			}
			pemStr := *getSecretResposne.Value
			err = os.WriteFile(tlsCertFilePath, []byte(pemStr), 0600)
			if err != nil {
				return nil, err
			}
			logger.Debug().Msgf("Stored certificate %s", p.config.TlsCertificateId)
		} else {
			return nil, err
		}
	}

	// pull jwk verification keys
	verifyJwks := make([]key.JsonWebKey, 0, len(p.config.JWTKeyCertIDs))
	for _, jwtCertID := range p.config.JWTKeyCertIDs {
		jwkFilePath := filepath.Join(p.versionedDir, fmt.Sprintf("jwk-%s.json", jwtCertID.ID()))
		if fileBytes, err := os.ReadFile(jwkFilePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// pull key from key vault
				resp, err := agentClient.GetCertificateWithResponse(c, jwtCertID.NamespaceKind(), jwtCertID.NamespaceID(), jwtCertID.ID())
				if err != nil {
					return nil, err
				}
				if resp.StatusCode() != http.StatusOK {
					return nil, fmt.Errorf("failed to get certificate: %d", resp.StatusCode())
				}
				logger.Debug().Any("resp", resp.JSON200).Any("status", resp.Status()).Msg("get certificate response")
				jwkBytes, err := json.Marshal(resp.JSON200.Jwk)
				if err != nil {
					return nil, err
				}
				if err := os.WriteFile(jwkFilePath, jwkBytes, 0600); err != nil {
					return nil, err
				}
				verifyJwks = append(verifyJwks, resp.JSON200.Jwk)
			} else {
				return nil, err
			}
		} else {
			jwk := key.JsonWebKey{}
			if err := json.Unmarshal(fileBytes, &jwk); err != nil {
				return nil, err
			}
			verifyJwks = append(verifyJwks, jwk)
		}
	}

	activeConfig := agentServerConfiguration{
		TLSCertificateFile: tlsCertFilePath,
		JWTVerifyKeys:      verifyJwks,
		fetchedConfig:      p.config,
	}

	if fileBytes, err := json.Marshal(activeConfig); err != nil {
		return nil, err
	} else if err := os.WriteFile(filepath.Join(p.versionedDir, "config.ready.json"), fileBytes, 0600); err != nil {
		return nil, err
	}

	return &activeConfig, nil
}

func (p *agentConfigServerProcessor) activeDirLink() string {
	return filepath.Join(p.configDir, "agent-server.active")
}

func (p *agentConfigServerProcessor) InitialLoad(c context.Context) (AgentServerConfiguration, error) {
	// load current active config
	logger := log.Ctx(c).With().Str("step", "initial load").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	if _, err := os.Lstat(p.activeDirLink()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	baseConfig := managedapp.AgentConfigServer{}
	readyConfig := agentServerConfiguration{}
	if err := utils.ReadJsonFile(filepath.Join(p.activeDirLink(), "config.json"), &baseConfig); err != nil {
		return nil, err
	}
	if err := utils.ReadJsonFile(filepath.Join(p.activeDirLink(), "config.ready.json"), &readyConfig); err != nil {
		return nil, err
	}
	readyConfig.fetchedConfig = &baseConfig

	p.readyConfig = &readyConfig
	return &readyConfig, nil
}

func (p *agentConfigServerProcessor) ProcessUpdate(c context.Context, nextConfig *managedapp.AgentConfigServer) (AgentServerConfiguration, bool, error) {
	logger := log.Ctx(c)
	if p.readyConfig != nil && p.readyConfig.Version() == nextConfig.Version {
		p.readyConfig.fetchedConfig = nextConfig
		return p.readyConfig, false, nil
	}
	p.configProvisioner = configProvisioner{
		processor:    p,
		versionedDir: filepath.Join(p.configDir, "versioned", fmt.Sprintf("agent-server.%s", nextConfig.Version)),
		config:       nextConfig,
	}
	nextReadyConfig, err := p.configProvisioner.provision(c)
	if err != nil {
		return nil, true, err
	}
	// make link
	linkName := p.activeDirLink()
	if _, err := os.Lstat(linkName); err == nil {
		if err := os.Remove(linkName); err != nil {
			logger.Error().Err(err).Msg("failed to remove symlink")
			return nil, true, err
		}
	}
	if err := os.Symlink(filepath.Join(".", "versioned", fmt.Sprintf("agent-server.%s", nextConfig.Version)), linkName); err != nil {
		logger.Error().Err(err).Msg("failed to create symlink")
		return nil, true, err
	}
	p.readyConfig = nextReadyConfig
	return nextReadyConfig, true, nil
}

func (p *agentConfigServerProcessor) Shutdown(c context.Context) error {
	if p.readyConfig.fetchedConfig != nil {

		toUpdate, err := json.Marshal(p.readyConfig.fetchedConfig)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(p.activeDirLink(), "config.json"), toUpdate, 0600)
		if err != nil {
			return err
		}
		log.Ctx(c).Debug().Msg("persisted config.json upon shutdown")
	}
	return nil
}

var _ AgentServerConfiguration = &agentServerConfiguration{}
