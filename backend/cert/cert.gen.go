// Package cert provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package cert

import (
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
	externalRef1 "github.com/stephenzsy/small-kms/backend/key"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for AzureKeyvaultResourceCategory.
const (
	AzureKeyvaultResourceCategoryCertificates AzureKeyvaultResourceCategory = "certificates"
	AzureKeyvaultResourceCategoryKeys         AzureKeyvaultResourceCategory = "keys"
	AzureKeyvaultResourceCategorySecrets      AzureKeyvaultResourceCategory = "secrets"
)

// Defines values for CertificateFlag.
const (
	CertificateFlagCA         CertificateFlag = "ca"
	CertificateFlagClientAuth CertificateFlag = "clientAuth"
	CertificateFlagRootCA     CertificateFlag = "rootCa"
	CertificateFlagServerAuth CertificateFlag = "serverAuth"
)

// AzureKeyvaultResourceCategory defines model for AzureKeyvaultResourceCategory.
type AzureKeyvaultResourceCategory string

// CertPolicy defines model for CertPolicy.
type CertPolicy = certPolicyComposed

// CertPolicyFields defines model for CertPolicyFields.
type CertPolicyFields struct {
	ExpiryTime                externalRef0.Period `json:"expiryTime"`
	Flags                     []CertificateFlag   `json:"flags"`
	IssuerNamespaceIdentifier externalRef0.Id     `json:"issuerNamespaceIdentifier"`
	// Deprecated:
	IssuerNamespaceKind     externalRef0.NamespaceKind   `json:"issuerNamespaceKind"`
	KeyExportable           bool                         `json:"keyExportable"`
	KeySpec                 externalRef1.SigningKeySpec  `json:"keySpec"`
	LifetimeAction          *externalRef1.LifetimeAction `json:"lifetimeAction,omitempty"`
	Subject                 CertificateSubject           `json:"subject"`
	SubjectAlternativeNames *SubjectAlternativeNames     `json:"subjectAlternativeNames,omitempty"`
	Version                 HexDigest                    `json:"version"`
}

// CertPolicyRef defines model for CertPolicyRef.
type CertPolicyRef = certPolicyRefComposed

// CertPolicyRefFields defines model for CertPolicyRefFields.
type CertPolicyRefFields struct {
	DisplayName string `json:"displayName"`
}

// Certificate defines model for Certificate.
type Certificate = certificateComposed

// CertificateAttributes defines model for CertificateAttributes.
type CertificateAttributes struct {
	Exp    *externalRef0.NumericDate     `json:"exp,omitempty"`
	Iat    *externalRef0.NumericDate     `json:"iat,omitempty"`
	Issuer *externalRef0.ResourceLocator `json:"issuer,omitempty"`
	Nbf    *externalRef0.NumericDate     `json:"nbf,omitempty"`
}

// CertificateFields defines model for CertificateFields.
type CertificateFields struct {
	Alg                     externalRef1.JsonWebSignatureAlgorithm `json:"alg"`
	Flags                   []CertificateFlag                      `json:"flags,omitempty"`
	Jwk                     externalRef1.JsonWebKey                `json:"jwk"`
	KeyVaultSecretID        string                                 `json:"sid,omitempty"`
	Subject                 CertificateSubject                     `json:"subject"`
	SubjectAlternativeNames *SubjectAlternativeNames               `json:"subjectAlternativeNames,omitempty"`
}

// CertificateFlag defines model for CertificateFlag.
type CertificateFlag string

// CertificateRef defines model for CertificateRef.
type CertificateRef = certificateRefComposed

// CertificateRefFields defines model for CertificateRefFields.
type CertificateRefFields struct {
	Attributes CertificateAttributes `json:"attributes"`
	Thumbprint string                `json:"thumbprint"`
}

// CertificateRuleIssuer defines model for CertificateRuleIssuer.
type CertificateRuleIssuer struct {
	CertificateId externalRef0.Id `json:"certificateId,omitempty"`
	PolicyId      externalRef0.Id `json:"policyId"`
}

// CertificateRuleMsEntraClientCredential defines model for CertificateRuleMsEntraClientCredential.
type CertificateRuleMsEntraClientCredential struct {
	CertificateIds []externalRef0.Id `json:"certificateIds,omitempty"`
	PolicyId       externalRef0.Id   `json:"policyId"`
}

// CertificateSubject defines model for CertificateSubject.
type CertificateSubject struct {
	CommonName string `json:"commonName"`
}

// SubjectAlternativeNames defines model for SubjectAlternativeNames.
type SubjectAlternativeNames struct {
	DNSNames    []string `json:"dnsNames,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	IPAddresses []net.IP `json:"ipAddresses,omitempty"`
}

// CertPolicyResponse defines model for CertPolicyResponse.
type CertPolicyResponse = CertPolicy

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = Certificate

// ListCertificatesParams defines parameters for ListCertificates.
type ListCertificatesParams struct {
	// PolicyId Policy ID
	PolicyId *string `form:"policyId,omitempty" json:"policyId,omitempty"`
}

// PutCertificateRuleIssuerJSONRequestBody defines body for PutCertificateRuleIssuer for application/json ContentType.
type PutCertificateRuleIssuerJSONRequestBody = CertificateRuleIssuer

// PutCertificateRuleMsEntraClientCredentialJSONRequestBody defines body for PutCertificateRuleMsEntraClientCredential for application/json ContentType.
type PutCertificateRuleMsEntraClientCredentialJSONRequestBody = CertificateRuleMsEntraClientCredential

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List cert policies
	// (GET /v1/{namespaceKind}/{namespaceId}/cert-policy)
	ListCertPolicies(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// Get cert policy
	// (GET /v1/{namespaceKind}/{namespaceId}/cert-policy/{resourceId})
	GetCertPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter) error
	// List Key Vault role assignments
	// (GET /v1/{namespaceKind}/{namespaceId}/cert-policy/{resourceId}/keyvault-role-assignments/{resourceCategory})
	ListKeyVaultRoleAssignments(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, resourceCategory AzureKeyvaultResourceCategory) error
	// Get certificate rules for namespace
	// (GET /v1/{namespaceKind}/{namespaceId}/cert-rule/issuer)
	GetCertificateRuleIssuer(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// Update certificate rules for namespace
	// (PUT /v1/{namespaceKind}/{namespaceId}/cert-rule/issuer)
	PutCertificateRuleIssuer(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// Get certificate rules for namespace
	// (GET /v1/{namespaceKind}/{namespaceId}/cert-rule/ms-entra-client-credential)
	GetCertificateRuleMsEntraClientCredential(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// Update certificate rules for namespace
	// (PUT /v1/{namespaceKind}/{namespaceId}/cert-rule/ms-entra-client-credential)
	PutCertificateRuleMsEntraClientCredential(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// List certificates
	// (GET /v1/{namespaceKind}/{namespaceId}/certificates)
	ListCertificates(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, params ListCertificatesParams) error
	// Get certificate
	// (GET /v1/{namespaceKind}/{namespaceId}/certificates/{resourceId})
	GetCertificate(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListCertPolicies converts echo context to params.
func (w *ServerInterfaceWrapper) ListCertPolicies(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListCertPolicies(ctx, namespaceKind, namespaceId)
	return err
}

// GetCertPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertPolicy(ctx echo.Context) error {
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

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertPolicy(ctx, namespaceKind, namespaceId, resourceId)
	return err
}

// ListKeyVaultRoleAssignments converts echo context to params.
func (w *ServerInterfaceWrapper) ListKeyVaultRoleAssignments(ctx echo.Context) error {
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

	// ------------- Path parameter "resourceCategory" -------------
	var resourceCategory AzureKeyvaultResourceCategory

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceCategory", runtime.ParamLocationPath, ctx.Param("resourceCategory"), &resourceCategory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceCategory: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeyVaultRoleAssignments(ctx, namespaceKind, namespaceId, resourceId, resourceCategory)
	return err
}

// GetCertificateRuleIssuer converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificateRuleIssuer(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificateRuleIssuer(ctx, namespaceKind, namespaceId)
	return err
}

// PutCertificateRuleIssuer converts echo context to params.
func (w *ServerInterfaceWrapper) PutCertificateRuleIssuer(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutCertificateRuleIssuer(ctx, namespaceKind, namespaceId)
	return err
}

// GetCertificateRuleMsEntraClientCredential converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificateRuleMsEntraClientCredential(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificateRuleMsEntraClientCredential(ctx, namespaceKind, namespaceId)
	return err
}

// PutCertificateRuleMsEntraClientCredential converts echo context to params.
func (w *ServerInterfaceWrapper) PutCertificateRuleMsEntraClientCredential(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutCertificateRuleMsEntraClientCredential(ctx, namespaceKind, namespaceId)
	return err
}

// ListCertificates converts echo context to params.
func (w *ServerInterfaceWrapper) ListCertificates(ctx echo.Context) error {
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

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListCertificatesParams
	// ------------- Optional query parameter "policyId" -------------

	err = runtime.BindQueryParameter("form", true, false, "policyId", ctx.QueryParams(), &params.PolicyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter policyId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListCertificates(ctx, namespaceKind, namespaceId, params)
	return err
}

// GetCertificate converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificate(ctx echo.Context) error {
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

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCertificate(ctx, namespaceKind, namespaceId, resourceId)
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

	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/cert-policy", wrapper.ListCertPolicies)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/cert-policy/:resourceId", wrapper.GetCertPolicy)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/cert-policy/:resourceId/keyvault-role-assignments/:resourceCategory", wrapper.ListKeyVaultRoleAssignments)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/cert-rule/issuer", wrapper.GetCertificateRuleIssuer)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceId/cert-rule/issuer", wrapper.PutCertificateRuleIssuer)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/cert-rule/ms-entra-client-credential", wrapper.GetCertificateRuleMsEntraClientCredential)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceId/cert-rule/ms-entra-client-credential", wrapper.PutCertificateRuleMsEntraClientCredential)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/certificates", wrapper.ListCertificates)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/certificates/:resourceId", wrapper.GetCertificate)

}
