// Package managedapp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package managedapp

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// ManagedApp defines model for ManagedApp.
type ManagedApp = ManagedAppRef

// ManagedAppParameters defines model for ManagedAppParameters.
type ManagedAppParameters struct {
	DisplayName                  string `json:"displayName"`
	SkipServicePrincipalCreation *bool  `json:"skipServicePrincipalCreation,omitempty"`
}

// ManagedAppRef defines model for ManagedAppRef.
type ManagedAppRef = managedAppRefComposed

// ManagedAppRefFields defines model for ManagedAppRefFields.
type ManagedAppRefFields struct {
	AppID openapi_types.UUID `json:"appId"`

	// ApplicationId Object ID
	ApplicationID      openapi_types.UUID `json:"applicationId"`
	ServicePrincipalID openapi_types.UUID `json:"servicePrincipalId"`
}

// CreateManagedAppJSONRequestBody defines body for CreateManagedApp for application/json ContentType.
type CreateManagedAppJSONRequestBody = ManagedAppParameters

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List managed apps
	// (GET /v1/managed-app)
	ListManagedApps(ctx echo.Context) error
	// Create a managed app
	// (POST /v1/managed-app)
	CreateManagedApp(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListManagedApps converts echo context to params.
func (w *ServerInterfaceWrapper) ListManagedApps(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListManagedApps(ctx)
	return err
}

// CreateManagedApp converts echo context to params.
func (w *ServerInterfaceWrapper) CreateManagedApp(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateManagedApp(ctx)
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

	router.GET(baseURL+"/v1/managed-app", wrapper.ListManagedApps)
	router.POST(baseURL+"/v1/managed-app", wrapper.CreateManagedApp)

}