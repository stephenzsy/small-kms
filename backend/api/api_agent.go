package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	agentconfig "github.com/stephenzsy/small-kms/backend/agent-config"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// AgentCheckIn implements models.ServerInterface.
func (*server) AgentCheckIn(ctx echo.Context, params models.AgentCheckInParams) error {
	return ctx.NoContent(http.StatusNoContent)
}

// AgentGetConfiguration implements models.ServerInterface.
func (*server) AgentGetConfiguration(ctx echo.Context, configName shared.AgentConfigName,
	params models.AgentGetConfigurationParams) error {
	bad := func(e error) error {
		return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)
	callerPrincipalId, err := auth.AuthorizeAgent(c)
	if err != nil {
		return bad(err)
	}

	c, err = ns.WithNamespaceContext(c, shared.NamespaceKindServicePrincipal, shared.UUIDIdentifier(callerPrincipalId))
	if err != nil {
		return bad(err)
	}
	config, err := agentconfig.GetAgentConfiguration(c, configName, &params)
	return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, config, err)
}

// GetAgentConfiguration implements models.ServerInterface.
func (*server) GetAgentConfiguration(ctx echo.Context, namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier, configName shared.AgentConfigName) error {
	bad := func(e error) error {
		return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	config, err := agentconfig.GetAgentConfiguration(c, configName, nil)
	return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, config, err)
}

// PutAgentConfiguration implements models.ServerInterface.
func (*server) PutAgentConfiguration(ctx echo.Context, namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier, configName shared.AgentConfigName) error {
	bad := func(e error) error {
		return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	configParams := shared.AgentConfigurationParameters{}
	err = c.Bind(&configParams)
	if err != nil {
		return bad(fmt.Errorf("%w:%w", common.ErrStatusBadRequest, err))
	}

	config, err := agentconfig.PutAgentConfiguration(c, configName, configParams)
	return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, config, err)
}
