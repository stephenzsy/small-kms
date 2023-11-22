package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

// PutAgentConfig implements admin.ServerInterface.
func (*AgentAdminServer) PutAgentConfig(ec echo.Context, namespaceId string, configName agentmodels.AgentConfigName) error {

	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	params := new(agentmodels.CreateAgentConfigRequest)
	if err := c.Bind(params); err != nil {
		return err
	}

	switch configName {
	case agentmodels.AgentConfigNameIdentity:
		return putAgentConfigIdentity(c, namespaceId, params)
	}
	return base.ErrResponseStatusNotFound

}