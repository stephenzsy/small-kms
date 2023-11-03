// Package agentpush provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package agentpush

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get agent diagnostics
	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/agent/instance/{resourceIdentifier}/diagnostics)
	GetAgentDiagnostics(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params GetAgentDiagnosticsParams) error

	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/agent/instance/{resourceIdentifier}/docker/images)
	AgentDockerImageList(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params AgentDockerImageListParams) error

	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/agent/instance/{resourceIdentifier}/docker/info)
	GetAgentDockerInfo(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params GetAgentDockerInfoParams) error

	// (POST /v1/{namespaceKind}/{namespaceIdentifier}/agent/instance/{resourceIdentifier}/pull-image)
	AgentPullImage(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params AgentPullImageParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAgentDiagnostics converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentDiagnostics(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAgentDiagnosticsParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Cryptocat-Proxy-Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Cryptocat-Proxy-Authorization")]; found {
		var XCryptocatProxyAuthorization DelegatedAuthorizationHeaderParameter
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Cryptocat-Proxy-Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, valueList[0], &XCryptocatProxyAuthorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Cryptocat-Proxy-Authorization: %s", err))
		}

		params.XCryptocatProxyAuthorization = &XCryptocatProxyAuthorization
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentDiagnostics(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	return err
}

// AgentDockerImageList converts echo context to params.
func (w *ServerInterfaceWrapper) AgentDockerImageList(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentDockerImageListParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Cryptocat-Proxy-Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Cryptocat-Proxy-Authorization")]; found {
		var XCryptocatProxyAuthorization DelegatedAuthorizationHeaderParameter
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Cryptocat-Proxy-Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, valueList[0], &XCryptocatProxyAuthorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Cryptocat-Proxy-Authorization: %s", err))
		}

		params.XCryptocatProxyAuthorization = &XCryptocatProxyAuthorization
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AgentDockerImageList(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	return err
}

// GetAgentDockerInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentDockerInfo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAgentDockerInfoParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Cryptocat-Proxy-Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Cryptocat-Proxy-Authorization")]; found {
		var XCryptocatProxyAuthorization DelegatedAuthorizationHeaderParameter
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Cryptocat-Proxy-Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, valueList[0], &XCryptocatProxyAuthorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Cryptocat-Proxy-Authorization: %s", err))
		}

		params.XCryptocatProxyAuthorization = &XCryptocatProxyAuthorization
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentDockerInfo(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	return err
}

// AgentPullImage converts echo context to params.
func (w *ServerInterfaceWrapper) AgentPullImage(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentPullImageParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Cryptocat-Proxy-Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Cryptocat-Proxy-Authorization")]; found {
		var XCryptocatProxyAuthorization DelegatedAuthorizationHeaderParameter
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Cryptocat-Proxy-Authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, valueList[0], &XCryptocatProxyAuthorization)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Cryptocat-Proxy-Authorization: %s", err))
		}

		params.XCryptocatProxyAuthorization = &XCryptocatProxyAuthorization
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AgentPullImage(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
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

	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/agent/instance/:resourceIdentifier/diagnostics", wrapper.GetAgentDiagnostics)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/agent/instance/:resourceIdentifier/docker/images", wrapper.AgentDockerImageList)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/agent/instance/:resourceIdentifier/docker/info", wrapper.GetAgentDockerInfo)
	router.POST(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/agent/instance/:resourceIdentifier/pull-image", wrapper.AgentPullImage)

}
