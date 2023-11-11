package radius

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
	if !ok {
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
	return pc, nil
}

func (p *radiusConfigProcessor) processClients(c context.Context) (string, error) {
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
		if err != nil {
			return "", err
		}
		clients[i].Secret = *kvResp.Value
	}
	marshalled, err := frconfig.FreeRadiusConfigList[frconfig.RadiusClientConfig](clients).MarshalFreeradiusConfig()
	if err != nil {
		return "", err
	}
	if err := configPath.WriteFile(marshalled); err != nil {
		return "", err
	}
	return configPath.Path(), nil
}
