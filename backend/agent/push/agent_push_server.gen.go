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

	// (POST /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId}/config/radius)
	PushAgentConfigRadius(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params PushAgentConfigRadiusParams) error

	// (DELETE /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId}/docker/containers/{containerId})
	AgentDockerContainerRemove(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, containerId string, params AgentDockerContainerRemoveParams) error

	// (GET /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId}/docker/containers/{containerId})
	AgentDockerContainerInspect(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, containerId string, params AgentDockerContainerInspectParams) error

	// (POST /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId}/docker/containers/{containerId}/stop)
	AgentDockerContainerStop(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, containerId string, params AgentDockerContainerStopParams) error

	// (POST /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId}/launch-agent)
	AgentLaunchAgent(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params AgentLaunchAgentParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PushAgentConfigRadius converts echo context to params.
func (w *ServerInterfaceWrapper) PushAgentConfigRadius(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PushAgentConfigRadiusParams

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
	err = w.Handler.PushAgentConfigRadius(ctx, namespaceKind, namespaceId, resourceId, params)
	return err
}

// AgentDockerContainerRemove converts echo context to params.
func (w *ServerInterfaceWrapper) AgentDockerContainerRemove(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	// ------------- Path parameter "containerId" -------------
	var containerId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "containerId", runtime.ParamLocationPath, ctx.Param("containerId"), &containerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter containerId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentDockerContainerRemoveParams

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
	err = w.Handler.AgentDockerContainerRemove(ctx, namespaceKind, namespaceId, resourceId, containerId, params)
	return err
}

// AgentDockerContainerInspect converts echo context to params.
func (w *ServerInterfaceWrapper) AgentDockerContainerInspect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	// ------------- Path parameter "containerId" -------------
	var containerId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "containerId", runtime.ParamLocationPath, ctx.Param("containerId"), &containerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter containerId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentDockerContainerInspectParams

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
	err = w.Handler.AgentDockerContainerInspect(ctx, namespaceKind, namespaceId, resourceId, containerId, params)
	return err
}

// AgentDockerContainerStop converts echo context to params.
func (w *ServerInterfaceWrapper) AgentDockerContainerStop(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	// ------------- Path parameter "containerId" -------------
	var containerId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "containerId", runtime.ParamLocationPath, ctx.Param("containerId"), &containerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter containerId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentDockerContainerStopParams

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
	err = w.Handler.AgentDockerContainerStop(ctx, namespaceKind, namespaceId, resourceId, containerId, params)
	return err
}

// AgentLaunchAgent converts echo context to params.
func (w *ServerInterfaceWrapper) AgentLaunchAgent(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params AgentLaunchAgentParams

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
	err = w.Handler.AgentLaunchAgent(ctx, namespaceKind, namespaceId, resourceId, params)
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

	router.POST(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId/config/radius", wrapper.PushAgentConfigRadius)
	router.DELETE(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId/docker/containers/:containerId", wrapper.AgentDockerContainerRemove)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId/docker/containers/:containerId", wrapper.AgentDockerContainerInspect)
	router.POST(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId/docker/containers/:containerId/stop", wrapper.AgentDockerContainerStop)
	router.POST(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId/launch-agent", wrapper.AgentLaunchAgent)

}
