package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

type AgentServerConfiguration struct {
}

type AgentConfigServerProcessor struct {
	envConfig         *agentutils.AgentEnv
	configDir         string
	readyVersion      string
	readyConfig       AgentServerConfiguration
	configProvisioner configProvisioner
}

type configProvisioner struct {
	processor    *AgentConfigServerProcessor
	versionedDir string
	config       *managedapp.AgentConfigServer
}

func (p *configProvisioner) provision(c context.Context) error {
	logger := log.Ctx(c)
	if _, err := os.Stat(p.versionedDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(p.versionedDir, 0700)
		} else {
			return err
		}
	}

	agentClient, err := p.processor.envConfig.AgentClient()
	if err != nil {
		return err
	}

	if fileBytes, err := json.Marshal(p.config); err != nil {
		return err
	} else if err := os.WriteFile(filepath.Join(p.versionedDir, "config.json"), fileBytes, 0600); err != nil {
		return err
	}

	tlsCertFilePath := filepath.Join(p.versionedDir, "tls-cert.pem")
	if _, err := os.Stat(tlsCertFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// pull certificate
			resp, err := agentClient.GetCertificateWithResponse(c, base.NamespaceKindServicePrincipal, base.StringIdentifier("me"), p.config.TlsCertificateId)
			if err != nil {
				return err
			}
			cert := resp.JSON200
			azSecretsClient, err := p.processor.envConfig.AzSecretsClient()
			if err != nil {
				return err
			}
			sid := azsecrets.ID(*cert.KeyVaultSecretID)
			getSecretResposne, err := azSecretsClient.GetSecret(c, sid.Name(), sid.Version(), nil)
			if err != nil {
				return err
			}
			pemStr := *getSecretResposne.Value
			err = os.WriteFile(tlsCertFilePath, []byte(pemStr), 0600)
			if err != nil {
				return err
			}
			logger.Debug().Msgf("Stored certificate %s", p.config.TlsCertificateId)
		} else {
			return err
		}
	}

	// pull jwk verification keys
	for _, jwtCertID := range p.config.JWTKeyCertIDs {
		jwkFilePath := filepath.Join(p.versionedDir, fmt.Sprintf("jwk-%s.json", jwtCertID.ResourceIdentifier()))
		if _, err := os.Stat(jwkFilePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// pull key from key vault
				resp, err := agentClient.GetCertificateWithResponse(c, jwtCertID.NamespaceKind(), jwtCertID.NamespaceIdentifier(), jwtCertID.ResourceIdentifier())
				if err != nil {
					return err
				}
				if resp.StatusCode() != http.StatusOK {
					return fmt.Errorf("failed to get certificate: %d", resp.StatusCode())
				}
				logger.Debug().Any("resp", resp.JSON200).Any("status", resp.Status()).Msg("get certificate response")
				jwkBytes, err := json.Marshal(resp.JSON200.Jwk)
				if err != nil {
					return err
				}
				if err := os.WriteFile(jwkFilePath, jwkBytes, 0600); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	logger.Debug().Any("config", p.config).Msg("provision agent config")

	return nil
}

func (p *AgentConfigServerProcessor) ProcessUpdate(c context.Context, nextConfig *managedapp.AgentConfigServer) error {
	if p.readyVersion == nextConfig.Version {
		// nothing to do, except update timestamp
		p.configProvisioner.config = nextConfig
		return nil
	}
	p.configProvisioner = configProvisioner{
		processor:    p,
		versionedDir: filepath.Join(p.configDir, "versioned", fmt.Sprintf("agent-server.%s", nextConfig.Version)),
		config:       nextConfig,
	}
	return p.configProvisioner.provision(c)
}

func (p *AgentConfigServerProcessor) Shutdown(c context.Context) error {
	// TODO: persist config json
	return nil
}
