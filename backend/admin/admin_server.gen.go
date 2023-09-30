// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get diagnostics
	// (GET /v1/diagnostics)
	GetDiagnosticsV1(c *gin.Context)
	// Get my profiles
	// (GET /v1/my/profiles)
	GetMyProfilesV1(c *gin.Context)
	// Sync my profiles
	// (POST /v1/my/profiles)
	SyncMyProfilesV1(c *gin.Context)
	// List namespaces
	// (GET /v1/namespaces/{namespaceType})
	ListNamespacesV1(c *gin.Context, namespaceType NamespaceType)
	// List policies
	// (GET /v1/{namespaceId}/policies)
	ListPoliciesV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Delete Certificate Policy
	// (DELETE /v1/{namespaceId}/policies/{policyIdentifier})
	DeletePolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyIdentifier string, params DeletePolicyV1Params)
	// Get Certificate Policy
	// (GET /v1/{namespaceId}/policies/{policyIdentifier})
	GetPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyIdentifier string)
	// Put Policy
	// (PUT /v1/{namespaceId}/policies/{policyIdentifier})
	PutPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyIdentifier string)
	// Apply policy
	// (POST /v1/{namespaceId}/policies/{policyId}/apply)
	ApplyPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyId openapi_types.UUID)
	// Get namespace profile
	// (GET /v1/{namespaceId}/profile)
	GetNamespaceProfileV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Register namespace
	// (POST /v1/{namespaceId}/profile)
	RegisterNamespaceProfileV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Link device service principal
	// (GET /v2/device/{namespaceId}/link-service-principal)
	GetDeviceServicePrincipalLinkV2(c *gin.Context, namespaceId NamespaceIdParameter, params GetDeviceServicePrincipalLinkV2Params)
	// Put certificate template
	// (POST /v2/group/{namespaceId}/certificate-templates/{templateId}/enroll)
	BeginEnrollCertificateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter, params BeginEnrollCertificateV2Params)
	// List namespaces by type
	// (GET /v2/{namespaceType})
	ListNamespacesByTypeV2(c *gin.Context, namespaceType NamespaceTypeParameter)
	// List certificate templates
	// (GET /v2/{namespaceType}/{namespaceId}/certificate-templates)
	ListCertificateTemplatesV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter)
	// Get certificate template
	// (GET /v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId})
	GetCertificateTemplateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// Put certificate template
	// (PUT /v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId})
	PutCertificateTemplateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// List certificates issued by template
	// (GET /v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}/certificates)
	ListCertificatesByTemplateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// Create certificate
	// (POST /v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}/certificates)
	IssueCertificateByTemplateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, templateId TemplateIdParameter, params IssueCertificateByTemplateV2Params)
	// Get certificate
	// (GET /v2/{namespaceType}/{namespaceId}/certificate-templates/{templateId}/certificates/latest)
	GetLatestCertificateByTemplateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, templateId TemplateIdParameter, params GetLatestCertificateByTemplateV2Params)
	// Get certificate
	// (GET /v2/{namespaceType}/{namespaceId}/certificates/{certId})
	GetCertificateV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, certId CertIdParameter, params GetCertificateV2Params)
	// complete certificate enrollment
	// (POST /v2/{namespaceType}/{namespaceId}/certificates/{certId}/pending)
	CompleteCertificateEnrollmentV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, certId CertIdParameter, params CompleteCertificateEnrollmentV2Params)
	// Sync namespace info with ms graph
	// (POST /v2/{namespaceType}/{namespaceId}/graph-sync)
	SyncNamespaceInfoV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetDiagnosticsV1 operation middleware
func (siw *ServerInterfaceWrapper) GetDiagnosticsV1(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetDiagnosticsV1(c)
}

// GetMyProfilesV1 operation middleware
func (siw *ServerInterfaceWrapper) GetMyProfilesV1(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetMyProfilesV1(c)
}

// SyncMyProfilesV1 operation middleware
func (siw *ServerInterfaceWrapper) SyncMyProfilesV1(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.SyncMyProfilesV1(c)
}

// ListNamespacesV1 operation middleware
func (siw *ServerInterfaceWrapper) ListNamespacesV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceType

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListNamespacesV1(c, namespaceType)
}

// ListPoliciesV1 operation middleware
func (siw *ServerInterfaceWrapper) ListPoliciesV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListPoliciesV1(c, namespaceId)
}

// DeletePolicyV1 operation middleware
func (siw *ServerInterfaceWrapper) DeletePolicyV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "policyIdentifier" -------------
	var policyIdentifier string

	err = runtime.BindStyledParameter("simple", false, "policyIdentifier", c.Param("policyIdentifier"), &policyIdentifier)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter policyIdentifier: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params DeletePolicyV1Params

	// ------------- Optional query parameter "purge" -------------

	err = runtime.BindQueryParameter("form", true, false, "purge", c.Request.URL.Query(), &params.Purge)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter purge: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeletePolicyV1(c, namespaceId, policyIdentifier, params)
}

// GetPolicyV1 operation middleware
func (siw *ServerInterfaceWrapper) GetPolicyV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "policyIdentifier" -------------
	var policyIdentifier string

	err = runtime.BindStyledParameter("simple", false, "policyIdentifier", c.Param("policyIdentifier"), &policyIdentifier)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter policyIdentifier: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetPolicyV1(c, namespaceId, policyIdentifier)
}

// PutPolicyV1 operation middleware
func (siw *ServerInterfaceWrapper) PutPolicyV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "policyIdentifier" -------------
	var policyIdentifier string

	err = runtime.BindStyledParameter("simple", false, "policyIdentifier", c.Param("policyIdentifier"), &policyIdentifier)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter policyIdentifier: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutPolicyV1(c, namespaceId, policyIdentifier)
}

// ApplyPolicyV1 operation middleware
func (siw *ServerInterfaceWrapper) ApplyPolicyV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "policyId" -------------
	var policyId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "policyId", c.Param("policyId"), &policyId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter policyId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ApplyPolicyV1(c, namespaceId, policyId)
}

// GetNamespaceProfileV1 operation middleware
func (siw *ServerInterfaceWrapper) GetNamespaceProfileV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetNamespaceProfileV1(c, namespaceId)
}

// RegisterNamespaceProfileV1 operation middleware
func (siw *ServerInterfaceWrapper) RegisterNamespaceProfileV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RegisterNamespaceProfileV1(c, namespaceId)
}

// GetDeviceServicePrincipalLinkV2 operation middleware
func (siw *ServerInterfaceWrapper) GetDeviceServicePrincipalLinkV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDeviceServicePrincipalLinkV2Params

	// ------------- Optional query parameter "apply" -------------

	err = runtime.BindQueryParameter("form", true, false, "apply", c.Request.URL.Query(), &params.Apply)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter apply: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetDeviceServicePrincipalLinkV2(c, namespaceId, params)
}

// BeginEnrollCertificateV2 operation middleware
func (siw *ServerInterfaceWrapper) BeginEnrollCertificateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params BeginEnrollCertificateV2Params

	// ------------- Optional query parameter "onBeHalfOf" -------------

	err = runtime.BindQueryParameter("form", true, false, "onBeHalfOf", c.Request.URL.Query(), &params.OnBeHalfOf)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter onBeHalfOf: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.BeginEnrollCertificateV2(c, namespaceId, templateId, params)
}

// ListNamespacesByTypeV2 operation middleware
func (siw *ServerInterfaceWrapper) ListNamespacesByTypeV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListNamespacesByTypeV2(c, namespaceType)
}

// ListCertificateTemplatesV2 operation middleware
func (siw *ServerInterfaceWrapper) ListCertificateTemplatesV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListCertificateTemplatesV2(c, namespaceType, namespaceId)
}

// GetCertificateTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetCertificateTemplateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCertificateTemplateV2(c, namespaceType, namespaceId, templateId)
}

// PutCertificateTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) PutCertificateTemplateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutCertificateTemplateV2(c, namespaceType, namespaceId, templateId)
}

// ListCertificatesByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) ListCertificatesByTemplateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListCertificatesByTemplateV2(c, namespaceType, namespaceId, templateId)
}

// IssueCertificateByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) IssueCertificateByTemplateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params IssueCertificateByTemplateV2Params

	// ------------- Optional query parameter "includeCertificate" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeCertificate", c.Request.URL.Query(), &params.IncludeCertificate)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter includeCertificate: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.IssueCertificateByTemplateV2(c, namespaceType, namespaceId, templateId, params)
}

// GetLatestCertificateByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetLatestCertificateByTemplateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "templateId" -------------
	var templateId TemplateIdParameter

	err = runtime.BindStyledParameter("simple", false, "templateId", c.Param("templateId"), &templateId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter templateId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLatestCertificateByTemplateV2Params

	// ------------- Optional query parameter "includeCertificate" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeCertificate", c.Request.URL.Query(), &params.IncludeCertificate)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter includeCertificate: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetLatestCertificateByTemplateV2(c, namespaceType, namespaceId, templateId, params)
}

// GetCertificateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetCertificateV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "certId" -------------
	var certId CertIdParameter

	err = runtime.BindStyledParameter("simple", false, "certId", c.Param("certId"), &certId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter certId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCertificateV2Params

	// ------------- Optional query parameter "includeCertificate" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeCertificate", c.Request.URL.Query(), &params.IncludeCertificate)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter includeCertificate: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCertificateV2(c, namespaceType, namespaceId, certId, params)
}

// CompleteCertificateEnrollmentV2 operation middleware
func (siw *ServerInterfaceWrapper) CompleteCertificateEnrollmentV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "certId" -------------
	var certId CertIdParameter

	err = runtime.BindStyledParameter("simple", false, "certId", c.Param("certId"), &certId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter certId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params CompleteCertificateEnrollmentV2Params

	// ------------- Optional query parameter "includeCertificate" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeCertificate", c.Request.URL.Query(), &params.IncludeCertificate)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter includeCertificate: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CompleteCertificateEnrollmentV2(c, namespaceType, namespaceId, certId, params)
}

// SyncNamespaceInfoV2 operation middleware
func (siw *ServerInterfaceWrapper) SyncNamespaceInfoV2(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceType" -------------
	var namespaceType NamespaceTypeParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceType", c.Param("namespaceType"), &namespaceType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceType: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId NamespaceIdParameter

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.SyncNamespaceInfoV2(c, namespaceType, namespaceId)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/v1/diagnostics", wrapper.GetDiagnosticsV1)
	router.GET(options.BaseURL+"/v1/my/profiles", wrapper.GetMyProfilesV1)
	router.POST(options.BaseURL+"/v1/my/profiles", wrapper.SyncMyProfilesV1)
	router.GET(options.BaseURL+"/v1/namespaces/:namespaceType", wrapper.ListNamespacesV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/policies", wrapper.ListPoliciesV1)
	router.DELETE(options.BaseURL+"/v1/:namespaceId/policies/:policyIdentifier", wrapper.DeletePolicyV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/policies/:policyIdentifier", wrapper.GetPolicyV1)
	router.PUT(options.BaseURL+"/v1/:namespaceId/policies/:policyIdentifier", wrapper.PutPolicyV1)
	router.POST(options.BaseURL+"/v1/:namespaceId/policies/:policyId/apply", wrapper.ApplyPolicyV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/profile", wrapper.GetNamespaceProfileV1)
	router.POST(options.BaseURL+"/v1/:namespaceId/profile", wrapper.RegisterNamespaceProfileV1)
	router.GET(options.BaseURL+"/v2/device/:namespaceId/link-service-principal", wrapper.GetDeviceServicePrincipalLinkV2)
	router.POST(options.BaseURL+"/v2/group/:namespaceId/certificate-templates/:templateId/enroll", wrapper.BeginEnrollCertificateV2)
	router.GET(options.BaseURL+"/v2/:namespaceType", wrapper.ListNamespacesByTypeV2)
	router.GET(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates", wrapper.ListCertificateTemplatesV2)
	router.GET(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates/:templateId", wrapper.GetCertificateTemplateV2)
	router.PUT(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates/:templateId", wrapper.PutCertificateTemplateV2)
	router.GET(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates/:templateId/certificates", wrapper.ListCertificatesByTemplateV2)
	router.POST(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates/:templateId/certificates", wrapper.IssueCertificateByTemplateV2)
	router.GET(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificate-templates/:templateId/certificates/latest", wrapper.GetLatestCertificateByTemplateV2)
	router.GET(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificates/:certId", wrapper.GetCertificateV2)
	router.POST(options.BaseURL+"/v2/:namespaceType/:namespaceId/certificates/:certId/pending", wrapper.CompleteCertificateEnrollmentV2)
	router.POST(options.BaseURL+"/v2/:namespaceType/:namespaceId/graph-sync", wrapper.SyncNamespaceInfoV2)
}
