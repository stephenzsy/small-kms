package radius

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
	frconfig "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
	"github.com/stephenzsy/small-kms/backend/base"
)

type radiusConfigProcessHandler struct {
	agentEnv  *agentcommon.AgentEnv
	configDir configmanager.ConfigDir
	hasLoaded bool
	processed ProcessedRadiusConfig
}

// After implements configmanager.ContextConfigHandler.
func (*radiusConfigProcessHandler) After(c context.Context) (context.Context, error) {
	return c, nil
}

// Before implements configmanager.ContextConfigHandler.
func (h *radiusConfigProcessHandler) Before(c context.Context) (context.Context, error) {

	config, ok := c.Value(contextKeyRadiusConfig).(*AgentConfigRadius)
	if !ok || config == nil {
		return c, nil
	}
	logger := log.Ctx(c)
	logger.Debug().Msg("radiusConfigProcessHandler.Before - begin")
	defer logger.Debug().Msg("radiusConfigProcessHandler.Before - done")

	if !h.hasLoaded {
		h.hasLoaded = true
		readyCfgPath := h.configDir.Active().File("config.ready.json")
		readyCfg, err := readyCfgPath.ReadFile()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				logger.Error().Err(err).Msgf("failed to read %s", readyCfgPath.Path())
			}
		} else if err := json.Unmarshal(readyCfg, &h.processed); err != nil {
			logger.Error().Err(err).Msgf("failed to unmarshal JSON: %s", readyCfgPath.Path())
		}
	}
	if h.processed.ConfigVersion == config.Version {
		h.processed.fetchedConfig = config
		return context.WithValue(c, contextKeyRadiusConfigProcessed, &h.processed), nil
	}

	processor := radiusConfigProcessor{
		config:    config,
		configDir: h.configDir,
		agentEnv:  h.agentEnv,
	}

	processed, err := processor.process(c)
	if err != nil {
		return c, err
	}
	h.processed = *processed
	h.processed.ConfigVersion = config.Version
	h.processed.fetchedConfig = config
	processedJson, err := json.Marshal(&h.processed)
	if err != nil {
		return c, err
	}
	if err := h.configDir.Versioned(config.Version).File("config.ready.json").WriteFile(processedJson); err != nil {
		return c, err
	}
	return context.WithValue(c, contextKeyRadiusConfigProcessed, &h.processed), nil

}

var _ configmanager.ContextConfigHandler = (*radiusConfigProcessHandler)(nil)

func NewRadiusConfigProcessHandler(agentEnv *agentcommon.AgentEnv, configDir configmanager.ConfigDir) *radiusConfigProcessHandler {
	return &radiusConfigProcessHandler{
		configDir: configDir,
		agentEnv:  agentEnv,
	}
}

type ProcessedRadiusConfig struct {
	ContainerVersion string   `json:"containerVersion"`
	ConfigVersion    string   `json:"configVersion"`
	HostBinds        []string `json:"hostBinds"`

	fetchedConfig *AgentConfigRadius
}

type radiusConfigProcessor struct {
	configDir configmanager.ConfigDir
	config    *AgentConfigRadius
	agentEnv  *agentcommon.AgentEnv

	agentClient agentclient.ClientWithResponsesInterface
}

func (p *radiusConfigProcessor) process(c context.Context) (*ProcessedRadiusConfig, error) {
	pc := &ProcessedRadiusConfig{}
	if p.agentClient == nil {
		if client, err := p.agentEnv.AgentClient(); err != nil {
			return nil, err
		} else {
			p.agentClient = client
		}
	}
	if configPath, err := p.processClients(c); err != nil {
		return nil, err
	} else {
		pc.HostBinds = append(pc.HostBinds, configPath+":/opt/etc/raddb/clients.conf:ro")
	}
	if hostBindings, err := p.processEAP(c); err != nil {
		return nil, err
	} else {
		pc.HostBinds = append(pc.HostBinds, hostBindings...)
	}
	if hostBinding, err := p.processServers(c); err != nil {
		return nil, err
	} else {
		pc.HostBinds = append(pc.HostBinds, hostBinding)
	}
	return pc, nil
}

func (p *radiusConfigProcessor) processClients(c context.Context) (string, error) {
	logger := log.Ctx(c)
	configPath := p.configDir.Versioned(p.config.Version).File("raddb", "clients.conf")
	if err := configPath.EnsureDirExist(); err != nil {
		return "", err
	}
	clients := p.config.Clients
	kvSClient, err := p.agentEnv.AzSecretsClient()
	if err != nil {
		return "", err
	}
	for i, client := range clients {
		if client.SecretId == "" {
			continue
		}
		// pull secret
		resp, err := p.agentClient.GetSecretWithResponse(c, base.NamespaceKindServicePrincipal, "me", client.SecretId, nil)
		if err != nil {
			return "", err
		}
		if resp.JSON200 == nil {
			return "", fmt.Errorf("no secret returned, status code %d", resp.StatusCode())
		}
		secretRef := azsecrets.ID(resp.JSON200.Sid)
		kvResp, err := kvSClient.GetSecret(c, secretRef.Name(), secretRef.Version(), nil)
		logger.Debug().Msgf("retrieved secret %s", secretRef)
		if err != nil {
			return "", err
		}
		clients[i].Secret = *kvResp.Value
	}
	sb := &strings.Builder{}
	if err := frconfig.FreeRadiusConfigList[frconfig.RadiusClientConfig](clients).MarshalFreeradiusConfig(sb, ""); err != nil {
		return "", err
	}
	if err := configPath.WriteFile([]byte(sb.String())); err != nil {
		return "", err
	}
	return configPath.Path(), nil
}

func (p *radiusConfigProcessor) processEAP(c context.Context) ([]string, error) {
	logger := log.Ctx(c)
	// fetch server cert

	kvSClient, err := p.agentEnv.AzSecretsClient()
	if err != nil {
		return nil, err
	}
	certResp, err := p.agentClient.GetCertificateWithResponse(c, base.NamespaceKindServicePrincipal, "me",
		p.config.EapTls.CertId)
	if err != nil {
		return nil, err
	}
	if certResp.JSON200 == nil {
		return nil, fmt.Errorf("no certificate returned, status code %d", certResp.StatusCode())
	}

	kvSecretID := azsecrets.ID(certResp.JSON200.KeyVaultSecretID)
	if kvSecretID == "" {
		return nil, fmt.Errorf("no keyvault secret id returned, status code %d", certResp.StatusCode())
	}
	kvResp, err := kvSClient.GetSecret(c, kvSecretID.Name(), kvSecretID.Version(), nil)
	if err != nil {
		return nil, err
	}
	logger.Info().Msgf("retrieved secret %s", kvSecretID)
	pemBytes := []byte(*kvResp.Value)
	raddbCertsDir := p.configDir.Versioned(p.config.Version).Dir("raddb", "certs")
	serverCertPath := raddbCertsDir.File("server.pem")
	if err := serverCertPath.EnsureDirExist(); err != nil {
		return nil, err
	}
	if err := serverCertPath.WriteFile(pemBytes); err != nil {
		return nil, err
	}

	// write CA file
	if len(certResp.JSON200.Jwk.CertificateChain) >= 2 {
		caCertPem := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte(certResp.JSON200.Jwk.CertificateChain[1]),
		})
		caCertPath := raddbCertsDir.File("ca.pem")
		if err := caCertPath.WriteFile(caCertPem); err != nil {
			return nil, err
		}
	}

	configPath := p.configDir.Versioned(p.config.Version).File("raddb", "mods-enabled", "eap")
	if err := configPath.EnsureDirExist(); err != nil {
		return nil, err
	}

	sb := &strings.Builder{}
	if err := frconfig.MarshalModEAP(sb, ""); err != nil {
		return nil, err
	}
	if err := configPath.WriteFile([]byte(sb.String())); err != nil {
		return nil, err
	}
	return []string{
		configPath.Path() + ":/opt/etc/raddb/mods-enabled/eap:ro",
		raddbCertsDir.Path() + ":/opt/etc/raddb/certs:ro",
	}, nil
}

func (p *radiusConfigProcessor) processServers(c context.Context) (string, error) {
	configDir := p.configDir.Versioned(p.config.Version).Dir("raddb", "sites-enabled")
	if err := configDir.EnsureDirExist(); err != nil {
		return "", err
	}
	for _, server := range p.config.Servers {

		configPath := configDir.File(server.Name)
		sb := &strings.Builder{}
		if err := server.MarshalFreeradiusConfig(sb, ""); err != nil {
			return "", err
		}
		if err := configPath.WriteFile([]byte(sb.String())); err != nil {
			return "", err
		}
	}
	return configDir.Path() + ":/opt/etc/raddb/sites-enabled:ro", nil

}
