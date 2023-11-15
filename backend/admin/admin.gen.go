// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List agents
	// (GET /v2/agents)
	ListAgents(ctx echo.Context) error
	// Get system app
	// (GET /v2/system-app/{id})
	GetSystemApp(ctx echo.Context, id string) error
	// Sync managed app
	// (POST /v2/system-app/{id})
	SyncSystemApp(ctx echo.Context, id string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListAgents converts echo context to params.
func (w *ServerInterfaceWrapper) ListAgents(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListAgents(ctx)
	return err
}

// GetSystemApp converts echo context to params.
func (w *ServerInterfaceWrapper) GetSystemApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetSystemApp(ctx, id)
	return err
}

// SyncSystemApp converts echo context to params.
func (w *ServerInterfaceWrapper) SyncSystemApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SyncSystemApp(ctx, id)
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

	router.GET(baseURL+"/v2/agents", wrapper.ListAgents)
	router.GET(baseURL+"/v2/system-app/:id", wrapper.GetSystemApp)
	router.POST(baseURL+"/v2/system-app/:id", wrapper.SyncSystemApp)

}
