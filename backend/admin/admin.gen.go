// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/models"
	externalRef1 "github.com/stephenzsy/small-kms/backend/models/agent"
	externalRef2 "github.com/stephenzsy/small-kms/backend/models/cert"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// ErrorResult defines model for ErrorResult.
type ErrorResult struct {
	Message *string `json:"message,omitempty"`
}

// IdParameter defines model for IdParameter.
type IdParameter = string

// NamespaceIdParameter defines model for NamespaceIdParameter.
type NamespaceIdParameter = string

// NamespaceProviderParameter defines model for NamespaceProviderParameter.
type NamespaceProviderParameter = externalRef0.NamespaceProvider

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = ErrorResult

// ListCertificatesParams defines parameters for ListCertificates.
type ListCertificatesParams struct {
	// PolicyId Policy ID
	PolicyId *string `form:"policyId,omitempty" json:"policyId,omitempty"`
}

// GetCertificateParams defines parameters for GetCertificate.
type GetCertificateParams struct {
	// IncludeJwk Include JWK
	IncludeJwk *bool `form:"includeJwk,omitempty" json:"includeJwk,omitempty"`
}

// CreateAgentJSONRequestBody defines body for CreateAgent for application/json ContentType.
type CreateAgentJSONRequestBody = externalRef1.CreateAgentRequest

// PutAgentConfigJSONRequestBody defines body for PutAgentConfig for application/json ContentType.
type PutAgentConfigJSONRequestBody = externalRef1.AgentConfigFields

// PutCertificatePolicyJSONRequestBody defines body for PutCertificatePolicy for application/json ContentType.
type PutCertificatePolicyJSONRequestBody = externalRef2.CreateCertificatePolicyRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create agent
	// (POST /v2/agents)
	CreateAgent(ctx echo.Context) error
	// Get agent
	// (GET /v2/agents/{id})
	GetAgent(ctx echo.Context, id IdParameter) error
	// list profiles
	// (GET /v2/profiles/{namespaceProvider})
	ListProfiles(ctx echo.Context, namespaceProvider NamespaceProviderParameter) error
	// Get agent config
	// (GET /v2/service-principal/{namespaceId}/agent-config)
	GetAgentConfig(ctx echo.Context, namespaceId NamespaceIdParameter) error
	// Put agent config
	// (PUT /v2/service-principal/{namespaceId}/agent-config)
	PutAgentConfig(ctx echo.Context, namespaceId NamespaceIdParameter) error
	// Get system app
	// (GET /v2/system-apps/{id})
	GetSystemApp(ctx echo.Context, id IdParameter) error
	// Sync managed app
	// (POST /v2/system-apps/{id})
	SyncSystemApp(ctx echo.Context, id IdParameter) error
	// List certificate policies
	// (GET /v2/{namespaceProvider}/{namespaceId}/certificate-policies)
	ListCertificatePolicies(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter) error
	// Get certificate policy
	// (GET /v2/{namespaceProvider}/{namespaceId}/certificate-policies/{id})
	GetCertificatePolicy(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
	// put certificate policy
	// (PUT /v2/{namespaceProvider}/{namespaceId}/certificate-policies/{id})
	PutCertificatePolicy(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
	// put certificate policy
	// (POST /v2/{namespaceProvider}/{namespaceId}/certificate-policies/{id}/generate)
	GenerateCertificate(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
	// List certificates
	// (GET /v2/{namespaceProvider}/{namespaceId}/certificates)
	ListCertificates(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, params ListCertificatesParams) error
	// Delete certificate
	// (DELETE /v2/{namespaceProvider}/{namespaceId}/certificates/{id})
	DeleteCertificate(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
	// Get certificate
	// (GET /v2/{namespaceProvider}/{namespaceId}/certificates/{id})
	GetCertificate(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter, params GetCertificateParams) error
	// List key policies
	// (GET /v2/{namespaceProvider}/{namespaceId}/key-policies)
	ListKeyPolicies(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter) error
	// Get key policy
	// (GET /v2/{namespaceProvider}/{namespaceId}/key-policies/{id})
	GetKeyPolicy(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
	// put key policy
	// (PUT /v2/{namespaceProvider}/{namespaceId}/key-policies/{id})
	PutKeyPolicy(ctx echo.Context, namespaceProvider NamespaceProviderParameter, namespaceId NamespaceIdParameter, id IdParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateAgent converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAgent(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateAgent(ctx)
	return err
}

// GetAgent converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgent(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgent(ctx, id)
	return err
}

// ListProfiles converts echo context to params.
func (w *ServerInterfaceWrapper) ListProfiles(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListProfiles(ctx, namespaceProvider)
	return err
}

// GetAgentConfig converts echo context to params.
func (w *ServerInterfaceWrapper) GetAgentConfig(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAgentConfig(ctx, namespaceId)
	return err
}

// PutAgentConfig converts echo context to params.
func (w *ServerInterfaceWrapper) PutAgentConfig(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutAgentConfig(ctx, namespaceId)
	return err
}

// GetSystemApp converts echo context to params.
func (w *ServerInterfaceWrapper) GetSystemApp(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id IdParameter

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
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SyncSystemApp(ctx, id)
	return err
}

// ListCertificatePolicies converts echo context to params.
func (w *ServerInterfaceWrapper) ListCertificatePolicies(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListCertificatePolicies(ctx, namespaceProvider, namespaceId)
	return err
}

// GetCertificatePolicy converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificatePolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificatePolicy(ctx, namespaceProvider, namespaceId, id)
	return err
}

// PutCertificatePolicy converts echo context to params.
func (w *ServerInterfaceWrapper) PutCertificatePolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutCertificatePolicy(ctx, namespaceProvider, namespaceId, id)
	return err
}

// GenerateCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) GenerateCertificate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GenerateCertificate(ctx, namespaceProvider, namespaceId, id)
	return err
}

// ListCertificates converts echo context to params.
func (w *ServerInterfaceWrapper) ListCertificates(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListCertificatesParams
	// ------------- Optional query parameter "policyId" -------------

	err = runtime.BindQueryParameter("form", true, false, "policyId", ctx.QueryParams(), &params.PolicyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter policyId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListCertificates(ctx, namespaceProvider, namespaceId, params)
	return err
}

// DeleteCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCertificate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteCertificate(ctx, namespaceProvider, namespaceId, id)
	return err
}

// GetCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCertificateParams
	// ------------- Optional query parameter "includeJwk" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeJwk", ctx.QueryParams(), &params.IncludeJwk)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter includeJwk: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificate(ctx, namespaceProvider, namespaceId, id, params)
	return err
}

// ListKeyPolicies converts echo context to params.
func (w *ServerInterfaceWrapper) ListKeyPolicies(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeyPolicies(ctx, namespaceProvider, namespaceId)
	return err
}

// GetKeyPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) GetKeyPolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetKeyPolicy(ctx, namespaceProvider, namespaceId, id)
	return err
}

// PutKeyPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) PutKeyPolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceProvider" -------------
	var namespaceProvider NamespaceProviderParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceProvider", runtime.ParamLocationPath, ctx.Param("namespaceProvider"), &namespaceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceProvider: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id IdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutKeyPolicy(ctx, namespaceProvider, namespaceId, id)
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

	router.POST(baseURL+"/v2/agents", wrapper.CreateAgent)
	router.GET(baseURL+"/v2/agents/:id", wrapper.GetAgent)
	router.GET(baseURL+"/v2/profiles/:namespaceProvider", wrapper.ListProfiles)
	router.GET(baseURL+"/v2/service-principal/:namespaceId/agent-config", wrapper.GetAgentConfig)
	router.PUT(baseURL+"/v2/service-principal/:namespaceId/agent-config", wrapper.PutAgentConfig)
	router.GET(baseURL+"/v2/system-apps/:id", wrapper.GetSystemApp)
	router.POST(baseURL+"/v2/system-apps/:id", wrapper.SyncSystemApp)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/certificate-policies", wrapper.ListCertificatePolicies)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/certificate-policies/:id", wrapper.GetCertificatePolicy)
	router.PUT(baseURL+"/v2/:namespaceProvider/:namespaceId/certificate-policies/:id", wrapper.PutCertificatePolicy)
	router.POST(baseURL+"/v2/:namespaceProvider/:namespaceId/certificate-policies/:id/generate", wrapper.GenerateCertificate)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/certificates", wrapper.ListCertificates)
	router.DELETE(baseURL+"/v2/:namespaceProvider/:namespaceId/certificates/:id", wrapper.DeleteCertificate)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/certificates/:id", wrapper.GetCertificate)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/key-policies", wrapper.ListKeyPolicies)
	router.GET(baseURL+"/v2/:namespaceProvider/:namespaceId/key-policies/:id", wrapper.GetKeyPolicy)
	router.PUT(baseURL+"/v2/:namespaceProvider/:namespaceId/key-policies/:id", wrapper.PutKeyPolicy)

}
