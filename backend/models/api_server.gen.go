// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
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
	// Provision agent
	// (POST /v3/application/{namespaceId}/agent)
	ProvisionAgentProfile(ctx echo.Context, namespaceId NamespaceIdParameter) error
	// Get diagnostics
	// (GET /v3/diagnostics)
	GetDiagnostics(ctx echo.Context) error
	// List profiles by type
	// (GET /v3/profiles/{namespaceKind})
	ListProfiles(ctx echo.Context, namespaceKind NamespaceKindParameter) error
	// Get namespace info with ms graph
	// (GET /v3/profiles/{namespaceKind}/{namespaceId})
	GetProfile(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter) error
	// Sync namespace info with ms graph
	// (POST /v3/profiles/{namespaceKind}/{namespaceId})
	SyncProfile(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter) error
	// Get service config
	// (GET /v3/service/config)
	GetServiceConfig(ctx echo.Context) error
	// Update service config
	// (PATCH /v3/service/config/{configPath})
	PatchServiceConfig(ctx echo.Context, configPath PatchServiceConfigParamsConfigPath) error

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
	// Delete certificate template
	// (DELETE /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId})
	DeleteCertificateTemplate(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter) error
	// List Key Vault role assignments
	// (GET /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments)
	ListKeyVaultRoleAssignments(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter) error
	// Add Key Vault role assignment
	// (POST /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments)
	AddKeyVaultRoleAssignment(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter, params AddKeyVaultRoleAssignmentParams) error
	// Remove Key Vault role assignment
	// (DELETE /v3/{namespaceKind}/{namespaceId}/certificate-template/{templateId}/keyvault-role-assignments/{roleAssignmentId})
	RemoveKeyVaultRoleAssignment(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, templateId CertificateTemplateIdentifierParameter, roleAssignmentId string) error
	// Create linked certificate template
	// (POST /v3/{namespaceKind}/{namespaceId}/certificate-templates)
	CreateLinkedCertificateTemplate(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter) error
	// Delete certificate
	// (DELETE /v3/{namespaceKind}/{namespaceId}/certificate/{certificateId})
	DeleteCertificate(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter) error
	// Get certificate
	// (GET /v3/{namespaceKind}/{namespaceId}/certificate/{certificateId})
	GetCertificate(ctx echo.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params GetCertificateParams) error
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

// ProvisionAgentProfile converts echo context to params.
func (w *ServerInterfaceWrapper) ProvisionAgentProfile(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ProvisionAgentProfile(ctx, namespaceId)
	return err
}

// GetDiagnostics converts echo context to params.
func (w *ServerInterfaceWrapper) GetDiagnostics(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDiagnostics(ctx)
	return err
}

// ListProfiles converts echo context to params.
func (w *ServerInterfaceWrapper) ListProfiles(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListProfiles(ctx, namespaceKind)
	return err
}

// GetProfile converts echo context to params.
func (w *ServerInterfaceWrapper) GetProfile(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProfile(ctx, namespaceKind, namespaceId)
	return err
}

// SyncProfile converts echo context to params.
func (w *ServerInterfaceWrapper) SyncProfile(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SyncProfile(ctx, namespaceKind, namespaceId)
	return err
}

// GetServiceConfig converts echo context to params.
func (w *ServerInterfaceWrapper) GetServiceConfig(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetServiceConfig(ctx)
	return err
}

// PatchServiceConfig converts echo context to params.
func (w *ServerInterfaceWrapper) PatchServiceConfig(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "configPath" -------------
	var configPath PatchServiceConfigParamsConfigPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "configPath", runtime.ParamLocationPath, ctx.Param("configPath"), &configPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter configPath: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PatchServiceConfig(ctx, configPath)
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

// DeleteCertificateTemplate converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCertificateTemplate(ctx echo.Context) error {
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
	err = w.Handler.DeleteCertificateTemplate(ctx, namespaceKind, namespaceId, templateId)
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

// CreateLinkedCertificateTemplate converts echo context to params.
func (w *ServerInterfaceWrapper) CreateLinkedCertificateTemplate(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateLinkedCertificateTemplate(ctx, namespaceKind, namespaceId)
	return err
}

// DeleteCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCertificate(ctx echo.Context) error {
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

	// ------------- Path parameter "certificateId" -------------
	var certificateId CertificateIdPathParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "certificateId", runtime.ParamLocationPath, ctx.Param("certificateId"), &certificateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter certificateId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteCertificate(ctx, namespaceKind, namespaceId, certificateId)
	return err
}

// GetCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificate(ctx echo.Context) error {
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

	// ------------- Path parameter "certificateId" -------------
	var certificateId CertificateIdPathParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "certificateId", runtime.ParamLocationPath, ctx.Param("certificateId"), &certificateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter certificateId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCertificateParams
	// ------------- Optional query parameter "includeCertificate" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeCertificate", ctx.QueryParams(), &params.IncludeCertificate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter includeCertificate: %s", err))
	}

	// ------------- Optional query parameter "templateId" -------------

	err = runtime.BindQueryParameter("form", true, false, "templateId", ctx.QueryParams(), &params.TemplateId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter templateId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificate(ctx, namespaceKind, namespaceId, certificateId, params)
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
	router.POST(baseURL+"/v3/application/:namespaceId/agent", wrapper.ProvisionAgentProfile)
	router.GET(baseURL+"/v3/diagnostics", wrapper.GetDiagnostics)
	router.GET(baseURL+"/v3/profiles/:namespaceKind", wrapper.ListProfiles)
	router.GET(baseURL+"/v3/profiles/:namespaceKind/:namespaceId", wrapper.GetProfile)
	router.POST(baseURL+"/v3/profiles/:namespaceKind/:namespaceId", wrapper.SyncProfile)
	router.GET(baseURL+"/v3/service/config", wrapper.GetServiceConfig)
	router.PATCH(baseURL+"/v3/service/config/:configPath", wrapper.PatchServiceConfig)
	router.GET(baseURL+"/v3/servicePrincipal/:namespaceId/agent-proxy", wrapper.GetAgentProxyInfo)
	router.GET(baseURL+"/v3/servicePrincipal/:namespaceId/agent-proxy/docker/info", wrapper.GetDockerInfo)
	router.POST(baseURL+"/v3/:namespaceKind/:namespaceId/agent-callback/:configName", wrapper.AgentCallback)
	router.GET(baseURL+"/v3/:namespaceKind/:namespaceId/agent-config/:configName", wrapper.GetAgentConfiguration)
	router.PUT(baseURL+"/v3/:namespaceKind/:namespaceId/agent-config/:configName", wrapper.PutAgentConfiguration)
	router.DELETE(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId", wrapper.DeleteCertificateTemplate)
	router.GET(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments", wrapper.ListKeyVaultRoleAssignments)
	router.POST(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments", wrapper.AddKeyVaultRoleAssignment)
	router.DELETE(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-template/:templateId/keyvault-role-assignments/:roleAssignmentId", wrapper.RemoveKeyVaultRoleAssignment)
	router.POST(baseURL+"/v3/:namespaceKind/:namespaceId/certificate-templates", wrapper.CreateLinkedCertificateTemplate)
	router.DELETE(baseURL+"/v3/:namespaceKind/:namespaceId/certificate/:certificateId", wrapper.DeleteCertificate)
	router.GET(baseURL+"/v3/:namespaceKind/:namespaceId/certificate/:certificateId", wrapper.GetCertificate)

}
