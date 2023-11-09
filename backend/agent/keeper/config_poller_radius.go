package keeper

import (
	"context"
	"fmt"
	"time"

	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

func NewRadiusConfigPoller(envConfig *agentcommon.AgentEnv) *configmanager.ConfigPoller[contextKey, *managedapp.AgentConfigRadius] {
	return configmanager.NewConfigPoller[contextKey, *managedapp.AgentConfigRadius]("radius-config-poller", contextKeyRadiusConfig, func(c context.Context) (*managedapp.AgentConfigRadius, time.Duration, error) {
		bad := func(err error) (*managedapp.AgentConfigRadius, time.Duration, error) {
			return nil, time.Minute * 5, err
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
		return resp.JSON200, time.Until(resp.JSON200.RefreshAfter), nil
	})
}
