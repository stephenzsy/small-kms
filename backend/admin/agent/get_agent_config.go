package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetAgentConfig implements admin.ServerInterface.
func (*AgentAdminServer) GetAgentConfig(ec echo.Context, namespaceId string, configName agentmodels.AgentConfigName) error {

	c := ec.(ctx.RequestContext)

	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	switch configName {
	case agentmodels.AgentConfigNameIdentity:
		return getAgentConfigIdentity(c, namespaceId)
	case agentmodels.AgentConfigNameEndpoint:
		return getAgentConfigEndpoint(c, namespaceId)
	}
	return base.ErrResponseStatusNotFound
}
