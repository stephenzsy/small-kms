package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

// GetAgentConfig implements admin.ServerInterface.
func (*AgentAdminServer) GetAgentConfig(ec echo.Context, namespaceId string, configName agentmodels.AgentConfigName) error {

	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	switch configName {
	case agentmodels.AgentConfigNameIdentity:
		return getAgentConfigIdentity(c, namespaceId)
	}
	return base.ErrResponseStatusNotFound
}
