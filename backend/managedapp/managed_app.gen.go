// Package managedapp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package managedapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
	externalRef0 "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
	externalRef1 "github.com/stephenzsy/small-kms/backend/base"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for AgentMode.
const (
	AgentModeLauncher AgentMode = "launcher"
	AgentModeServer   AgentMode = "server"
)

// AgentConfig defines model for AgentConfig.
type AgentConfig = agentConfigComposed

// AgentConfigFields defines model for AgentConfigFields.
type AgentConfigFields struct {
	RefreshAfter time.Time `json:"refreshAfter"`
	Version      string    `json:"version"`
}

// AgentConfigRadius defines model for AgentConfigRadius.
type AgentConfigRadius = agentConfigRadiusComposed

// AgentConfigRadiusFields defines model for AgentConfigRadiusFields.
type AgentConfigRadiusFields struct {
	Clients   []externalRef0.RadiusClientConfig `json:"clients,omitempty"`
	Container *AgentContainerConfiguration      `json:"container,omitempty"`
	DebugMode *bool                             `json:"debugMode,omitempty"`
	EapTls    *externalRef0.RadiusEapTls        `json:"eapTls,omitempty"`
	Servers   []externalRef0.RadiusServerConfig `json:"servers,omitempty"`
}

// AgentConfigServer defines model for AgentConfigServer.
type AgentConfigServer = agentConfigServerComposed

// AgentConfigServerEnv Environment variables for the agent config server, must be set manually
type AgentConfigServerEnv struct {
	EnvVarAzureContainerRegistryImageRepository string `json:"AZURE_ACR_IMAGE_REPOSITORY"`
	EnvVarAzureKeyVaultResourceEndpoint         string `json:"AZURE_KEYVAULT_RESOURCEENDPOINT"`
	Message                                     string `json:"_message"`
}

// AgentConfigServerFields defines model for AgentConfigServerFields.
type AgentConfigServerFields struct {
	AzureACRImageRef string `json:"azureAcrImageRef"`

	// Env Environment variables for the agent config server, must be set manually
	Env                    AgentConfigServerEnv           `json:"env"`
	JWTKeyCertIDs          []externalRef1.ResourceLocator `json:"jwtKeyCertIds"`
	JwtKeyCertPolicyId     externalRef1.ResourceLocator   `json:"jwtKeyCertPolicyId"`
	TlsCertificateId       externalRef1.Id                `json:"tlsCertificateId"`
	TlsCertificatePolicyId externalRef1.Id                `json:"tlsCertificatePolicyId"`
}

// AgentContainerConfiguration defines model for AgentContainerConfiguration.
type AgentContainerConfiguration struct {
	ContainerName    string        `json:"containerName,omitempty"`
	Env              []string      `json:"env,omitempty"`
	ExposedPortSpecs []string      `json:"exposedPortSpecs"`
	HostBinds        []string      `json:"hostBinds"`
	ImageRepo        string        `json:"imageRepo"`
	ImageTag         string        `json:"imageTag"`
	NetworkName      string        `json:"networkName,omitempty"`
	Secrets          []SecretMount `json:"secrets,omitempty"`
}

// AgentInstance defines model for AgentInstance.
type AgentInstance = agentInstanceComposed

// AgentInstanceFields defines model for AgentInstanceFields.
type AgentInstanceFields struct {
	BuildID  string    `json:"buildId"`
	Endpoint string    `json:"endpoint,omitempty"`
	Mode     AgentMode `json:"mode"`
	Version  string    `json:"version"`
}

// AgentMode defines model for AgentMode.
type AgentMode string

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
	ServicePrincipalType string             `json:"servicePrincipalType,omitempty"`
}

// SecretMount defines model for SecretMount.
type SecretMount struct {
	Source     string `json:"source"`
	TargetName string `json:"targetName"`
}

// ManagedAppIdParameter defines model for ManagedAppIdParameter.
type ManagedAppIdParameter = openapi_types.UUID

// AgentConfigRadiusResponse defines model for AgentConfigRadiusResponse.
type AgentConfigRadiusResponse = AgentConfigRadius

// CreateManagedAppJSONRequestBody defines body for CreateManagedApp for application/json ContentType.
type CreateManagedAppJSONRequestBody = ManagedAppParameters

// PatchAgentConfigRadiusJSONRequestBody defines body for PatchAgentConfigRadius for application/json ContentType.
type PatchAgentConfigRadiusJSONRequestBody = AgentConfigRadiusFields

// PutAgentConfigServerJSONRequestBody defines body for PutAgentConfigServer for application/json ContentType.
type PutAgentConfigServerJSONRequestBody = AgentConfigServerFields

// PutAgentInstanceJSONRequestBody defines body for PutAgentInstance for application/json ContentType.
type PutAgentInstanceJSONRequestBody = AgentInstanceFields

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
	// Get agent config radius
	// (GET /v1/{namespaceKind}/{namespaceId}/agent-config/radius)
	GetAgentConfigRadius(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter) error
	// Patch agent config radius
	// (PATCH /v1/{namespaceKind}/{namespaceId}/agent-config/radius)
	PatchAgentConfigRadius(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter) error
	// Get agent config server
	// (GET /v1/{namespaceKind}/{namespaceId}/agent-config/server)
	GetAgentConfigServer(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter) error
	// Put agent config server
	// (PUT /v1/{namespaceKind}/{namespaceId}/agent-config/server)
	PutAgentConfigServer(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter) error
	// List agent config server instances
	// (GET /v1/{namespaceKind}/{namespaceId}/agent/instance)
	ListAgentInstances(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter) error
	// Get agent config server instance
	// (GET /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId})
	GetAgentInstance(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter, resourceId externalRef1.ResourceIdParameter) error
	// Put agent config server instance
	// (PUT /v1/{namespaceKind}/{namespaceId}/agent/instance/{resourceId})
	PutAgentInstance(ctx echo.Context, namespaceKind externalRef1.NamespaceKindParameter, namespaceId externalRef1.NamespaceIdParameter, resourceId externalRef1.ResourceIdParameter) error
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

// GetAgentConfigRadius converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentConfigRadius(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentConfigRadius(ctx, namespaceKind, namespaceId)
	return err
}

// PatchAgentConfigRadius converts echo context to params.
func (w *ServerInterfaceWrapper) PatchAgentConfigRadius(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PatchAgentConfigRadius(ctx, namespaceKind, namespaceId)
	return err
}

// GetAgentConfigServer converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentConfigServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentConfigServer(ctx, namespaceKind, namespaceId)
	return err
}

// PutAgentConfigServer converts echo context to params.
func (w *ServerInterfaceWrapper) PutAgentConfigServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutAgentConfigServer(ctx, namespaceKind, namespaceId)
	return err
}

// ListAgentInstances converts echo context to params.
func (w *ServerInterfaceWrapper) ListAgentInstances(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListAgentInstances(ctx, namespaceKind, namespaceId)
	return err
}

// GetAgentInstance converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentInstance(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef1.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentInstance(ctx, namespaceKind, namespaceId, resourceId)
	return err
}

// PutAgentInstance converts echo context to params.
func (w *ServerInterfaceWrapper) PutAgentInstance(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef1.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef1.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef1.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutAgentInstance(ctx, namespaceKind, namespaceId, resourceId)
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
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/agent-config/radius", wrapper.GetAgentConfigRadius)
	router.PATCH(baseURL+"/v1/:namespaceKind/:namespaceId/agent-config/radius", wrapper.PatchAgentConfigRadius)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/agent-config/server", wrapper.GetAgentConfigServer)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceId/agent-config/server", wrapper.PutAgentConfigServer)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance", wrapper.ListAgentInstances)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId", wrapper.GetAgentInstance)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceId/agent/instance/:resourceId", wrapper.PutAgentInstance)

}
