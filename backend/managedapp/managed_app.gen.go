// Package managedapp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package managedapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for SystemAppName.
const (
	SystemAppNameAPI     SystemAppName = "api"
	SystemAppNameBackend SystemAppName = "backend"
)

// AgentConfig defines model for AgentConfig.
type AgentConfig = agentConfigComposed

// AgentConfigFields defines model for AgentConfigFields.
type AgentConfigFields struct {
	RefreshAfter time.Time `json:"refreshAfter"`
	Version      string    `json:"version"`
}

// AgentConfigServer defines model for AgentConfigServer.
type AgentConfigServer = agentConfigServerComposed

// AgentConfigServerFields defines model for AgentConfigServerFields.
type AgentConfigServerFields struct {
	ImageTag         string                                  `json:"imageTag"`
	JWTKeyCertIDs    []externalRef0.ResourceUniqueIdentifier `json:"jwtKeyCertIds"`
	TlsCertificateId externalRef0.Identifier                 `json:"tlsCertificateId"`
}

// AgentConfigServerParameters defines model for AgentConfigServerParameters.
type AgentConfigServerParameters struct {
	JwtKeyCertPolicyId     externalRef0.ResourceUniqueIdentifier `json:"jwtKeyCertPolicyId"`
	TlsCertificatePolicyId externalRef0.Identifier               `json:"tlsCertificatePolicyId"`
}

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
	ApplicationID        openapi_types.UUID `json:"applicationId"`
	ServicePrincipalID   openapi_types.UUID `json:"servicePrincipalId"`
	ServicePrincipalType *string            `json:"servicePrincipalType,omitempty"`
}

// SystemAppName defines model for SystemAppName.
type SystemAppName string

// ManagedAppIdParameter defines model for ManagedAppIdParameter.
type ManagedAppIdParameter = openapi_types.UUID

// CreateManagedAppJSONRequestBody defines body for CreateManagedApp for application/json ContentType.
type CreateManagedAppJSONRequestBody = ManagedAppParameters

// PutAgentConfigServerJSONRequestBody defines body for PutAgentConfigServer for application/json ContentType.
type PutAgentConfigServerJSONRequestBody = AgentConfigServerParameters

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List managed apps
	// (GET /v1/managed-app)
	ListManagedApps(ctx echo.Context) error
	// Create a managed app
	// (POST /v1/managed-app)
	CreateManagedApp(ctx echo.Context) error
	// Get managed app
	// (GET /v1/managed-app/{managedAppId})
	GetManagedApp(ctx echo.Context, managedAppId ManagedAppIdParameter) error
	// Sync managed app
	// (POST /v1/managed-app/{managedAppId})
	SyncManagedApp(ctx echo.Context, managedAppId ManagedAppIdParameter) error
	// Get system app
	// (GET /v1/system-app/{systemAppName})
	GetSystemApp(ctx echo.Context, systemAppName SystemAppName) error
	// Sync managed app
	// (POST /v1/system-app/{systemAppName})
	SyncSystemApp(ctx echo.Context, systemAppName SystemAppName) error
	// Put agent config server
	// (PUT /v1/{namespaceKind}/{namespaceIdentifier}/agent-config/server)
	PutAgentConfigServer(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter) error
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

// GetManagedApp converts echo context to params.
func (w *ServerInterfaceWrapper) GetManagedApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "managedAppId" -------------
	var managedAppId ManagedAppIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "managedAppId", runtime.ParamLocationPath, ctx.Param("managedAppId"), &managedAppId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter managedAppId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetManagedApp(ctx, managedAppId)
	return err
}

// SyncManagedApp converts echo context to params.
func (w *ServerInterfaceWrapper) SyncManagedApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "managedAppId" -------------
	var managedAppId ManagedAppIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "managedAppId", runtime.ParamLocationPath, ctx.Param("managedAppId"), &managedAppId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter managedAppId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SyncManagedApp(ctx, managedAppId)
	return err
}

// GetSystemApp converts echo context to params.
func (w *ServerInterfaceWrapper) GetSystemApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "systemAppName" -------------
	var systemAppName SystemAppName

	err = runtime.BindStyledParameterWithLocation("simple", false, "systemAppName", runtime.ParamLocationPath, ctx.Param("systemAppName"), &systemAppName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter systemAppName: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetSystemApp(ctx, systemAppName)
	return err
}

// SyncSystemApp converts echo context to params.
func (w *ServerInterfaceWrapper) SyncSystemApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "systemAppName" -------------
	var systemAppName SystemAppName

	err = runtime.BindStyledParameterWithLocation("simple", false, "systemAppName", runtime.ParamLocationPath, ctx.Param("systemAppName"), &systemAppName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter systemAppName: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SyncSystemApp(ctx, systemAppName)
	return err
}

// PutAgentConfigServer converts echo context to params.
func (w *ServerInterfaceWrapper) PutAgentConfigServer(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutAgentConfigServer(ctx, namespaceKind, namespaceIdentifier)
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
	router.GET(baseURL+"/v1/managed-app/:managedAppId", wrapper.GetManagedApp)
	router.POST(baseURL+"/v1/managed-app/:managedAppId", wrapper.SyncManagedApp)
	router.GET(baseURL+"/v1/system-app/:systemAppName", wrapper.GetSystemApp)
	router.POST(baseURL+"/v1/system-app/:systemAppName", wrapper.SyncSystemApp)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/agent-config/server", wrapper.PutAgentConfigServer)

}
