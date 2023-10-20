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
func (*server) AgentCallback(
	ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	configName shared.AgentConfigName) error {
	if configName != shared.AgentConfigNameActiveServer {
		return ctx.NoContent(http.StatusNoContent)
	}
	c := ctx.(RequestContext)
	namespaceId, err := ns.ResolveAuthedNamespaseID(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	req := shared.AgentConfiguration{}
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("%w:%w", common.ErrStatusBadRequest, err)
	}
	return wrapEchoResponse(c, agentconfig.ApiRecordAgentActiveServerCallback(c, &req))
}

// GetAgentConfiguration implements models.ServerInterface.
func (*server) GetAgentConfiguration(ctx echo.Context, namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier, configName shared.AgentConfigName,
	params models.GetAgentConfigurationParams) error {
	bad := func(e error) error {
		return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)
	var isAdmin bool
	if auth.AuthorizeAdminOnly(c) == nil {
		isAdmin = true
	}
	namespaceId, err := ns.ResolveAuthedNamespaseID(c, namespaceKind, namespaceId)
	if err != nil && !isAdmin {
		return bad(err)
	}

	c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	config, err := agentconfig.GetAgentConfiguration(c, configName, &params, isAdmin)
	return wrapResponse[*models.AgentConfigurationResponse](ctx, http.StatusOK, config, err)
}

// PutAgentConfiguration implements models.ServerInterface.
func (*server) PutAgentConfiguration(
	ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	configName shared.AgentConfigName) error {
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

// GetAgentProfile implements models.ServerInterface.
func (*server) GetAgentProfile(ctx echo.Context, namespaceId shared.Identifier) error {
	c := ctx.(RequestContext)
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return wrapEchoResponse(c, err)
	}
	c, err := ns.WithNamespaceContext(c, shared.NamespaceKindApplication, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	return wrapEchoResponse(c, agentconfig.ApiGetAgentProfile(c))
}

// ProvisionAgentProfile implements models.ServerInterface.
func (*server) ProvisionAgentProfile(ctx echo.Context, namespaceId shared.Identifier) error {
	c := ctx.(RequestContext)
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return wrapEchoResponse(c, err)
	}
	c, err := ns.WithNamespaceContext(c, shared.NamespaceKindApplication, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	params := shared.AgentProfileParameters{}
	if err := c.Bind(&params); err != nil {
		return wrapEchoResponse(c, fmt.Errorf("%w:%w", common.ErrStatusBadRequest, err))
	}
	return wrapEchoResponse(c, agentconfig.ApiProvisionAgentProfile(c, &params))
}
