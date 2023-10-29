// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Provision agent
	// (GET /v3/application/{namespaceId}/agent)
	GetAgentProfile(ctx echo.Context, namespaceId NamespaceIdParameter) error

	// (GET /v3/servicePrincipal/{namespaceId}/agent-proxy)
	GetAgentProxyInfo(ctx echo.Context, namespaceId NamespaceIdParameter) error

	// (GET /v3/servicePrincipal/{namespaceId}/agent-proxy/docker/info)
	GetDockerInfo(ctx echo.Context, namespaceId NamespaceIdParameter) error

	// (POST /v3/{namespaceKind}/{namespaceId}/agent-callback/{configName})
	AgentCallback(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter) error
	// Get agent autoconfig
	// (GET /v3/{namespaceKind}/{namespaceId}/agent-config/{configName})
	GetAgentConfiguration(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params GetAgentConfigurationParams) error
	// Get agent autoconfig
	// (PUT /v3/{namespaceKind}/{namespaceId}/agent-config/{configName})
	PutAgentConfiguration(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter) error
	// List Key Vault role assignments
	// (GET /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments)
	ListKeyVaultRoleAssignments(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter) error
	// Add Key Vault role assignment
	// (POST /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments)
	AddKeyVaultRoleAssignment(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter, params AddKeyVaultRoleAssignmentParams) error
	// Remove Key Vault role assignment
	// (DELETE /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments/{roleAssignmentId})
	RemoveKeyVaultRoleAssignment(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter, roleAssignmentId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAgentProfile converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentProfile(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentProfile(ctx, namespaceId)
	return err
}

// GetAgentProxyInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentProxyInfo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentProxyInfo(ctx, namespaceId)
	return err
}

// GetDockerInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetDockerInfo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDockerInfo(ctx, namespaceId)
	return err
}

// AgentCallback converts echo context to params.
func (w *ServerInterfaceWrapper) AgentCallback(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "configName" -------------
	var configName AgentConfigNameParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "configName", runtime.ParamLocationPath, ctx.Param("configName"), &configName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter configName: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AgentCallback(ctx, namespaceKind, namespaceId, configName)
	return err
}

// GetAgentConfiguration converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentConfiguration(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "configName" -------------
	var configName AgentConfigNameParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "configName", runtime.ParamLocationPath, ctx.Param("configName"), &configName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter configName: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAgentConfigurationParams
	// ------------- Optional query parameter "refreshToken" -------------

	err = runtime.BindQueryParameter("form", true, false, "refreshToken", ctx.QueryParams(), &params.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter refreshToken: %s", err))
	}

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Smallkms-If-Version-Not-Match" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Smallkms-If-Version-Not-Match")]; found {
		var XSmallkmsIfVersionNotMatch string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Smallkms-If-Version-Not-Match, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Smallkms-If-Version-Not-Match", runtime.ParamLocationHeader, valueList[0], &XSmallkmsIfVersionNotMatch)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Smallkms-If-Version-Not-Match: %s", err))
		}

		params.XSmallkmsIfVersionNotMatch = &XSmallkmsIfVersionNotMatch
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentConfiguration(ctx, namespaceKind, namespaceId, configName, params)
	return err
}

// PutAgentConfiguration converts echo context to params.
func (w *ServerInterfaceWrapper) PutAgentConfiguration(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "configName" -------------
	var configName AgentConfigNameParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "configName", runtime.ParamLocationPath, ctx.Param("configName"), &configName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter configName: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutAgentConfiguration(ctx, namespaceKind, namespaceId, configName)
	return err
}

// ListKeyVaultRoleAssignments converts echo context to params.
func (w *ServerInterfaceWrapper) ListKeyVaultRoleAssignments(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "templateId" -------------
	var templateId CertificateTemplateIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "templateId", runtime.ParamLocationPath, ctx.Param("templateId"), &templateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter templateId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeyVaultRoleAssignments(ctx, namespaceKind, namespaceId, templateId)
	return err
}

// AddKeyVaultRoleAssignment converts echo context to params.
func (w *ServerInterfaceWrapper) AddKeyVaultRoleAssignment(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "templateId" -------------
	var templateId CertificateTemplateIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "templateId", runtime.ParamLocationPath, ctx.Param("templateId"), &templateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter templateId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AddKeyVaultRoleAssignmentParams
	// ------------- Required query parameter "roleDefinitionId" -------------

	err = runtime.BindQueryParameter("form", true, true, "roleDefinitionId", ctx.QueryParams(), &params.RoleDefinitionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roleDefinitionId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AddKeyVaultRoleAssignment(ctx, namespaceKind, namespaceId, templateId, params)
	return err
}

// RemoveKeyVaultRoleAssignment converts echo context to params.
func (w *ServerInterfaceWrapper) RemoveKeyVaultRoleAssignment(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "templateId" -------------
	var templateId CertificateTemplateIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "templateId", runtime.ParamLocationPath, ctx.Param("templateId"), &templateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter templateId: %s", err))
	}

	// ------------- Path parameter "roleAssignmentId" -------------
	var roleAssignmentId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "roleAssignmentId", runtime.ParamLocationPath, ctx.Param("roleAssignmentId"), &roleAssignmentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roleAssignmentId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RemoveKeyVaultRoleAssignment(ctx, namespaceKind, namespaceId, templateId, roleAssignmentId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v3/application/:namespaceId/agent", wrapper.GetAgentProfile)
	router.GET(baseURL+"/v3/servicePrincipal/:namespaceId/agent-proxy", wrapper.GetAgentProxyInfo)
	router.GET(baseURL+"/v3/servicePrincipal/:namespaceId/agent-proxy/docker/info", wrapper.GetDockerInfo)
	router.POST(baseURL+"/v3/:namespaceKind/:namespaceId/agent-callback/:configName", wrapper.AgentCallback)
	router.GET(baseURL+"/v3/:namespaceKind/:namespaceId/agent-config/:configName", wrapper.GetAgentConfiguration)
	router.PUT(baseURL+"/v3/:namespaceKind/:namespaceId/agent-config/:configName", wrapper.PutAgentConfiguration)
	router.GET(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments", wrapper.ListKeyVaultRoleAssignments)
	router.POST(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments", wrapper.AddKeyVaultRoleAssignment)
	router.DELETE(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments/:roleAssignmentId", wrapper.RemoveKeyVaultRoleAssignment)

}
