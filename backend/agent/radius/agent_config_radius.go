package radius

import (
	"context"
	"fmt"

	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

type RadiusConfigManager struct {
	configmanager.ConfigPoller[contextKey, managedapp.AgentConfigRadius]

	envConfig *agentcommon.AgentEnv
}

func NewRadiusConfigManager(handlerChain *configmanager.ChainedContextConfigHandler, envConfig *agentcommon.AgentEnv, configDir string) *RadiusConfigManager {
	return &RadiusConfigManager{
		ConfigPoller: *configmanager.NewConfigPoller[contextKey, managedapp.AgentConfigRadius](
			handlerChain,
			"radius-config-poller",
			contextKeyRadiusConfig, func(c context.Context) (*managedapp.AgentConfigRadius, error) {
				bad := func(err error) (*managedapp.AgentConfigRadius, error) {
					return nil, err
				}
				client, err := envConfig.AgentClient()
				if err != nil {
					return bad(err)
				}
				resp, err := client.GetAgentConfigRadiusWithResponse(c, base.NamespaceKindServicePrincipal, base.ID("me"))
				if err != nil {
					return bad(err)
				}
				if resp.JSON200 == nil {
					return bad(fmt.Errorf("no radius config, status code %d", resp.StatusCode()))
				}
				return resp.JSON200, nil
			},
			configmanager.NewConfigCache[managedapp.AgentConfigRadius]("radius", configDir)),
		envConfig: envConfig,
	}
}
