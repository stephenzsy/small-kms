// Package agentserver provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package agentserver

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/shared"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AgentConfigNameParameter defines model for AgentConfigNameParameter.
type AgentConfigNameParameter = externalRef0.AgentConfigName

// CertificateTemplateIdentifierParameter defines model for CertificateTemplateIdentifierParameter.
type CertificateTemplateIdentifierParameter = externalRef0.Identifier

// NamespaceIdParameter defines model for NamespaceIdParameter.
type NamespaceIdParameter = externalRef0.Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = externalRef0.NamespaceKind

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v3/servicePrincipal/{namespaceId}/agent-proxy/docker/info)
	GetDockerInfo(ctx echo.Context, namespaceId NamespaceIdParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
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

	router.GET(baseURL+"/v3/servicePrincipal/:namespaceId/agent-proxy/docker/info", wrapper.GetDockerInfo)

}
