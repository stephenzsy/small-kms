package radius

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
	frconfig "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
	"github.com/stephenzsy/small-kms/backend/base"
)

type radiusConfigProcessHandler struct {
	agentEnv  *agentcommon.AgentEnv
	configDir configmanager.ConfigDir
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

	processor := radiusConfigProcessor{
		config:    config,
		configDir: h.configDir,
		agentEnv:  h.agentEnv,
	}

	return processor.process(c)
}

var _ configmanager.ContextConfigHandler = (*radiusConfigProcessHandler)(nil)

func NewRadiusConfigProcessHandler(agentEnv *agentcommon.AgentEnv, configDir configmanager.ConfigDir) *radiusConfigProcessHandler {
	return &radiusConfigProcessHandler{
		configDir: configDir,
		agentEnv:  agentEnv,
	}
}

type ProcessedRadiusConfig struct {
}

type radiusConfigProcessor struct {
	configDir configmanager.ConfigDir
	config    *AgentConfigRadius
	agentEnv  *agentcommon.AgentEnv

	agentClient agentclient.ClientWithResponsesInterface
}

func (p *radiusConfigProcessor) process(c context.Context) (context.Context, error) {
	if p.agentClient == nil {
		if client, err := p.agentEnv.AgentClient(); err != nil {
			return c, err
		} else {
			p.agentClient = client
		}
	}
	if err := p.processClients(c); err != nil {
		return c, err
	}
	return c, nil
}

func (p *radiusConfigProcessor) processClients(c context.Context) error {
	configPath := p.configDir.Versioned(p.config.Version).File("raddb", "clients.conf")
	if err := configPath.EnsureDirExist(); err != nil {
		return err
	}
	clients := p.config.Clients
	kvSClient, err := p.agentEnv.AzSecretsClient()
	if err != nil {
		return err
	}
	for i, client := range clients {
		if client.SecretId == "" {
			continue
		}
		// pull secret
		resp, err := p.agentClient.GetSecretWithResponse(c, base.NamespaceKindServicePrincipal, "me", client.SecretId, nil)
		if err != nil {
			return err
		}
		if resp.JSON200 == nil {
			return fmt.Errorf("no secret returned, status code %d", resp.StatusCode())
		}
		secretRef := azsecrets.ID(resp.JSON200.Sid)
		kvResp, err := kvSClient.GetSecret(c, secretRef.Name(), secretRef.Version(), nil)
		if err != nil {
			return err
		}
		clients[i].Secret = *kvResp.Value
	}
	marshalled, err := frconfig.FreeRadiusConfigList[frconfig.RadiusClientConfig](clients).MarshalFreeradiusConfig()
	if err != nil {
		return err
	}
	if err := configPath.WriteFile(marshalled); err != nil {
		return err
	}
	return nil
}
