// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get diagnostics
	// (GET /v1/diagnostics)
	GetDiagnosticsV1(c *gin.Context)
	// List certificate templates
	// (GET /v2/{namespaceId}/certificate-templates)
	ListCertificateTemplatesV2(c *gin.Context, namespaceId NamespaceIdParameter, params ListCertificateTemplatesV2Params)
	// Get certificate template
	// (GET /v2/{namespaceId}/certificate-templates/{templateId})
	GetCertificateTemplateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// Put certificate template
	// (PUT /v2/{namespaceId}/certificate-templates/{templateId})
	PutCertificateTemplateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// List certificates issued by template
	// (GET /v2/{namespaceId}/certificate-templates/{templateId}/certificates)
	ListCertificatesByTemplateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// Create certificate
	// (POST /v2/{namespaceId}/certificate-templates/{templateId}/certificates)
	IssueCertificateByTemplateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter, params IssueCertificateByTemplateV2Params)
	// Get certificate
	// (GET /v2/{namespaceId}/certificate-templates/{templateId}/certificates/latest)
	GetLatestCertificateByTemplateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter, params GetLatestCertificateByTemplateV2Params)
	// Put certificate template
	// (POST /v2/{namespaceId}/certificate-templates/{templateId}/enroll)
	BeginEnrollCertificateV2(c *gin.Context, namespaceId NamespaceIdParameter, templateId TemplateIdParameter)
	// Get certificate
	// (GET /v2/{namespaceId}/certificates/{certId})
	GetCertificateV2(c *gin.Context, namespaceId NamespaceIdParameter, certId CertIdParameter, params GetCertificateV2Params)
	// complete certificate enrollment
	// (POST /v2/{namespaceId}/certificates/{certId}/pending)
	CompleteCertificateEnrollmentV2(c *gin.Context, namespaceId NamespaceIdParameter, certId CertIdParameter, params CompleteCertificateEnrollmentV2Params)
	// Link device service principal
	// (GET /v2/{namespaceId}/link-service-principal)
	GetDeviceServicePrincipalLinkV2(c *gin.Context, namespaceId NamespaceIdParameter)
	// Link device service principal
	// (POST /v2/{namespaceId}/link-service-principal)
	CreateDeviceServicePrincipalLinkV2(c *gin.Context, namespaceId NamespaceIdParameter)
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

// ListCertificateTemplatesV2 operation middleware
func (siw *ServerInterfaceWrapper) ListCertificateTemplatesV2(c *gin.Context) {

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
	var params ListCertificateTemplatesV2Params

	// ------------- Optional query parameter "includeDefaultForType" -------------

	err = runtime.BindQueryParameter("form", true, false, "includeDefaultForType", c.Request.URL.Query(), &params.IncludeDefaultForType)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter includeDefaultForType: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListCertificateTemplatesV2(c, namespaceId, params)
}

// GetCertificateTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetCertificateTemplateV2(c *gin.Context) {

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

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCertificateTemplateV2(c, namespaceId, templateId)
}

// PutCertificateTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) PutCertificateTemplateV2(c *gin.Context) {

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

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutCertificateTemplateV2(c, namespaceId, templateId)
}

// ListCertificatesByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) ListCertificatesByTemplateV2(c *gin.Context) {

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

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListCertificatesByTemplateV2(c, namespaceId, templateId)
}

// IssueCertificateByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) IssueCertificateByTemplateV2(c *gin.Context) {

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

	siw.Handler.IssueCertificateByTemplateV2(c, namespaceId, templateId, params)
}

// GetLatestCertificateByTemplateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetLatestCertificateByTemplateV2(c *gin.Context) {

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

	siw.Handler.GetLatestCertificateByTemplateV2(c, namespaceId, templateId, params)
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

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.BeginEnrollCertificateV2(c, namespaceId, templateId)
}

// GetCertificateV2 operation middleware
func (siw *ServerInterfaceWrapper) GetCertificateV2(c *gin.Context) {

	var err error

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

	siw.Handler.GetCertificateV2(c, namespaceId, certId, params)
}

// CompleteCertificateEnrollmentV2 operation middleware
func (siw *ServerInterfaceWrapper) CompleteCertificateEnrollmentV2(c *gin.Context) {

	var err error

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

	siw.Handler.CompleteCertificateEnrollmentV2(c, namespaceId, certId, params)
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

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetDeviceServicePrincipalLinkV2(c, namespaceId)
}

// CreateDeviceServicePrincipalLinkV2 operation middleware
func (siw *ServerInterfaceWrapper) CreateDeviceServicePrincipalLinkV2(c *gin.Context) {

	var err error

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

	siw.Handler.CreateDeviceServicePrincipalLinkV2(c, namespaceId)
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
	router.GET(options.BaseURL+"/v2/:namespaceId/certificate-templates", wrapper.ListCertificateTemplatesV2)
	router.GET(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId", wrapper.GetCertificateTemplateV2)
	router.PUT(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId", wrapper.PutCertificateTemplateV2)
	router.GET(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId/certificates", wrapper.ListCertificatesByTemplateV2)
	router.POST(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId/certificates", wrapper.IssueCertificateByTemplateV2)
	router.GET(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId/certificates/latest", wrapper.GetLatestCertificateByTemplateV2)
	router.POST(options.BaseURL+"/v2/:namespaceId/certificate-templates/:templateId/enroll", wrapper.BeginEnrollCertificateV2)
	router.GET(options.BaseURL+"/v2/:namespaceId/certificates/:certId", wrapper.GetCertificateV2)
	router.POST(options.BaseURL+"/v2/:namespaceId/certificates/:certId/pending", wrapper.CompleteCertificateEnrollmentV2)
	router.GET(options.BaseURL+"/v2/:namespaceId/link-service-principal", wrapper.GetDeviceServicePrincipalLinkV2)
	router.POST(options.BaseURL+"/v2/:namespaceId/link-service-principal", wrapper.CreateDeviceServicePrincipalLinkV2)
}
