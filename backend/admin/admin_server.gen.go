// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.14.0 DO NOT EDIT.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for CertificateUsage.
const (
	UsageClientOnly      CertificateUsage = "client-only"
	UsageIntCA           CertificateUsage = "intermediate-ca"
	UsageRootCA          CertificateUsage = "root-ca"
	UsageServerAndClient CertificateUsage = "server-and-client"
	UsageServerOnly      CertificateUsage = "server-only"
)

// Defines values for KeyPropertiesCrv.
const (
	EcCurveP256 KeyPropertiesCrv = "P-256"
	EcCurveP384 KeyPropertiesCrv = "P-384"
)

// Defines values for KeyPropertiesKeySize.
const (
	KeySize2048 KeyPropertiesKeySize = 2048
	KeySize3072 KeyPropertiesKeySize = 3072
	KeySize4096 KeyPropertiesKeySize = 4096
)

// Defines values for KeyPropertiesKty.
const (
	KtyEC  KeyPropertiesKty = "EC"
	KtyRSA KeyPropertiesKty = "RSA"
)

// Defines values for NamespaceType.
const (
	NamespaceTypeBuiltInCaInt            NamespaceType = "#builtin.ca.intermediate"
	NamespaceTypeMsGraphDevice           NamespaceType = "#microsoft.graph.device"
	NamespaceTypeMsGraphGroup            NamespaceType = "#microsoft.graph.group"
	NamespaceTypeMsGraphServicePrincipal NamespaceType = "#microsoft.graph.servicePrincipal"
	NamespaceTypeMsGraphUser             NamespaceType = "#microsoft.graph.user"
)

// Defines values for PolicyType.
const (
	PolicyTypeCertRequest PolicyType = "certRequest"
)

// Defines values for GetCertificateV1ParamsAccept.
const (
	AcceptJson       GetCertificateV1ParamsAccept = "application/json"
	AcceptPem        GetCertificateV1ParamsAccept = "application/x-pem-file"
	AcceptX509CaCert GetCertificateV1ParamsAccept = "application/x-x509-ca-cert"
)

// ApplyPolicyRequest defines model for ApplyPolicyRequest.
type ApplyPolicyRequest struct {
	// CheckConsistency Check consistency of the policy
	CheckConsistency *bool `json:"checkConsistency,omitempty"`

	// ForceRenewCertificate Force certificate renewal
	ForceRenewCertificate *bool `json:"forceRenewCertificate,omitempty"`
}

// CertificateLifetimeTrigger defines model for CertificateLifetimeTrigger.
type CertificateLifetimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

// CertificateRef defines model for CertificateRef.
type CertificateRef struct {
	// CreatedBy Unique ID of the user who created the certificate
	CreatedBy string `json:"createdBy"`

	// Id Unique ID of the certificate, also the serial number of the certificate
	ID openapi_types.UUID `json:"id"`

	// Issuer Issuer of the certificate
	Issuer openapi_types.UUID `json:"issuer"`

	// IssuerNamespace Issuer of the certificate
	IssuerNamespace openapi_types.UUID `json:"issuerNamespace"`

	// Name Name of the certificate, also the common name (CN) of the certificate
	Name string `json:"name"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`

	// NotAfter Expiration date of the certificate
	NotAfter time.Time        `json:"notAfter"`
	Usage    CertificateUsage `json:"usage"`
}

// CertificateRequestPolicyParameters defines model for CertificateRequestPolicyParameters.
type CertificateRequestPolicyParameters struct {
	// IssuerNamespaceId ID of the issuer namespace
	IssuerNamespaceID        openapi_types.UUID                  `json:"issuerNamespaceId"`
	KeyProperties            *KeyProperties                      `json:"keyProperties,omitempty"`
	KeyStorePath             string                              `json:"keyStorePath"`
	LifetimeTrigger          *CertificateLifetimeTrigger         `json:"lifetimeTrigger,omitempty"`
	MsGraphGroupAllowMembers *bool                               `json:"msGraphGroupAllowMembers,omitempty"`
	Subject                  CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames  *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                    CertificateUsage                    `json:"usage"`
	ValidityInMonths         *int32                              `json:"validity_months,omitempty"`
}

// CertificateSubject defines model for CertificateSubject.
type CertificateSubject struct {
	// C Country or region
	C *string `json:"c,omitempty"`

	// Cn Common name
	CN string `json:"cn"`

	// O Organization
	O *string `json:"o,omitempty"`

	// Ou Organizational unit
	OU *string `json:"ou,omitempty"`
}

// CertificateSubjectAlternativeNames defines model for CertificateSubjectAlternativeNames.
type CertificateSubjectAlternativeNames struct {
	DNSNames           *[]string `json:"dns_names,omitempty"`
	Emails             *[]string `json:"emails,omitempty"`
	UserPrincipalNames *[]string `json:"upns,omitempty"`
}

// CertificateUsage defines model for CertificateUsage.
type CertificateUsage string

// KeyProperties defines model for KeyProperties.
type KeyProperties struct {
	CurveName *KeyPropertiesCrv     `json:"crv,omitempty"`
	KeySize   *KeyPropertiesKeySize `json:"key_size,omitempty"`
	KeyType   KeyPropertiesKty      `json:"kty"`

	// ReuseKey Keep using the same key version if exists
	ReuseKey *bool `json:"reuse_key,omitempty"`
}

// KeyPropertiesCrv defines model for KeyProperties.Crv.
type KeyPropertiesCrv string

// KeyPropertiesKeySize defines model for KeyProperties.KeySize.
type KeyPropertiesKeySize int32

// KeyPropertiesKty defines model for KeyProperties.Kty.
type KeyPropertiesKty string

// NamespaceProfile defines model for NamespaceProfile.
type NamespaceProfile struct {
	// DeviceId \#microsoft.graph.device deviceId
	DeviceID        *openapi_types.UUID `json:"deviceId,omitempty"`
	DeviceOwnership *string             `json:"deviceOwnership,omitempty"`
	DisplayName     string              `json:"displayName"`
	ID              openapi_types.UUID  `json:"id"`

	// IsCompliant \#microsoft.graph.device isCompliant
	IsCompliant *bool `json:"isCompliant,omitempty"`

	// Manufacturer \#microsoft.graph.device manufacturer
	Manufacturer *string         `json:"manufacturer,omitempty"`
	MemberOf     *[]NamespaceRef `json:"memberOf,omitempty"`

	// Model \#microsoft.graph.device model
	Model *string `json:"model,omitempty"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	ObjectType  NamespaceType      `json:"objectType"`

	// OperatingSystem \#microsoft.graph.device operatingSystem
	OperatingSystem *string `json:"operatingSystem,omitempty"`

	// OperatingSystemVersion \#microsoft.graph.device operatingSystemVersion
	OperatingSystemVersion *string         `json:"operatingSystemVersion,omitempty"`
	RegisterdOwners        *[]NamespaceRef `json:"registerdOwners,omitempty"`
	ServicePrincipalType   *string         `json:"servicePrincipalType,omitempty"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy         string  `json:"updatedBy"`
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
}

// NamespaceRef defines model for NamespaceRef.
type NamespaceRef struct {
	DisplayName string             `json:"displayName"`
	ID          openapi_types.UUID `json:"id"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	ObjectType  NamespaceType      `json:"objectType"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// NamespaceType defines model for NamespaceType.
type NamespaceType string

// Policy defines model for Policy.
type Policy struct {
	CertRequest *CertificateRequestPolicyParameters `json:"certRequest,omitempty"`
	ID          openapi_types.UUID                  `json:"id"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	PolicyType  PolicyType         `json:"policyType"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// PolicyParameters defines model for PolicyParameters.
type PolicyParameters struct {
	CertRequest *CertificateRequestPolicyParameters `json:"certRequest,omitempty"`
	PolicyType  PolicyType                          `json:"policyType"`
}

// PolicyRef defines model for PolicyRef.
type PolicyRef struct {
	ID openapi_types.UUID `json:"id"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	PolicyType  PolicyType         `json:"policyType"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// PolicyState defines model for PolicyState.
type PolicyState struct {
	CertRequest *PolicyStateCertRequest `json:"certRequest,omitempty"`
	ID          openapi_types.UUID      `json:"id"`
	Message     string                  `json:"message"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	PolicyType  PolicyType         `json:"policyType"`
	Status      *string            `json:"status,omitempty"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// PolicyStateCertRequest defines model for PolicyStateCertRequest.
type PolicyStateCertRequest struct {
	LastAction      string             `json:"lastAction"`
	LastCertExpires time.Time          `json:"lastCertExpires"`
	LastCertID      openapi_types.UUID `json:"lastCertId"`
	LastCertIssued  time.Time          `json:"lastCertIssued"`
}

// PolicyType defines model for PolicyType.
type PolicyType string

// RequestDiagnostics defines model for RequestDiagnostics.
type RequestDiagnostics struct {
	RequestHeaders []RequestHeaderEntry              `json:"requestHeaders"`
	ServiceRuntime RequestDiagnostics_ServiceRuntime `json:"serviceRuntime"`
}

// RequestDiagnostics_ServiceRuntime defines model for RequestDiagnostics.ServiceRuntime.
type RequestDiagnostics_ServiceRuntime struct {
	GoVersion            string            `json:"goVersion"`
	AdditionalProperties map[string]string `json:"-"`
}

// RequestHeaderEntry defines model for RequestHeaderEntry.
type RequestHeaderEntry struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// ResourceRef defines model for ResourceRef.
type ResourceRef struct {
	ID openapi_types.UUID `json:"id"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Code                 *string                `json:"code,omitempty"`
	Message              *string                `json:"message,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// GetCertificateV1Params defines parameters for GetCertificateV1.
type GetCertificateV1Params struct {
	Accept *GetCertificateV1ParamsAccept `json:"Accept,omitempty"`
}

// GetCertificateV1ParamsAccept defines parameters for GetCertificateV1.
type GetCertificateV1ParamsAccept string

// PutPolicyV1JSONRequestBody defines body for PutPolicyV1 for application/json ContentType.
type PutPolicyV1JSONRequestBody = PolicyParameters

// ApplyPolicyV1JSONRequestBody defines body for ApplyPolicyV1 for application/json ContentType.
type ApplyPolicyV1JSONRequestBody = ApplyPolicyRequest

// Getter for additional properties for RequestDiagnostics_ServiceRuntime. Returns the specified
// element and whether it was found
func (a RequestDiagnostics_ServiceRuntime) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for RequestDiagnostics_ServiceRuntime
func (a *RequestDiagnostics_ServiceRuntime) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for RequestDiagnostics_ServiceRuntime to handle AdditionalProperties
func (a *RequestDiagnostics_ServiceRuntime) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["goVersion"]; found {
		err = json.Unmarshal(raw, &a.GoVersion)
		if err != nil {
			return fmt.Errorf("error reading 'goVersion': %w", err)
		}
		delete(object, "goVersion")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for RequestDiagnostics_ServiceRuntime to handle AdditionalProperties
func (a RequestDiagnostics_ServiceRuntime) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["goVersion"], err = json.Marshal(a.GoVersion)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'goVersion': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ErrorResponse. Returns the specified
// element and whether it was found
func (a ErrorResponse) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ErrorResponse
func (a *ErrorResponse) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ErrorResponse to handle AdditionalProperties
func (a *ErrorResponse) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["code"]; found {
		err = json.Unmarshal(raw, &a.Code)
		if err != nil {
			return fmt.Errorf("error reading 'code': %w", err)
		}
		delete(object, "code")
	}

	if raw, found := object["message"]; found {
		err = json.Unmarshal(raw, &a.Message)
		if err != nil {
			return fmt.Errorf("error reading 'message': %w", err)
		}
		delete(object, "message")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ErrorResponse to handle AdditionalProperties
func (a ErrorResponse) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	if a.Code != nil {
		object["code"], err = json.Marshal(a.Code)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'code': %w", err)
		}
	}

	if a.Message != nil {
		object["message"], err = json.Marshal(a.Message)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'message': %w", err)
		}
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get diagnostics
	// (GET /v1/diagnostics)
	GetDiagnosticsV1(c *gin.Context)
	// Get my profiles
	// (GET /v1/my/profiles)
	GetMyProfilesV1(c *gin.Context)
	// List namespaces
	// (GET /v1/namespaces/{namespaceType})
	ListNamespacesV1(c *gin.Context, namespaceType NamespaceType)
	// List certificates
	// (GET /v1/{namespaceId}/certificates)
	ListCertificatesV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Get certificate
	// (GET /v1/{namespaceId}/certificates/{id})
	GetCertificateV1(c *gin.Context, namespaceId openapi_types.UUID, id openapi_types.UUID, params GetCertificateV1Params)
	// List policies
	// (GET /v1/{namespaceId}/policies)
	ListPoliciesV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Get Certificate Policy
	// (GET /v1/{namespaceId}/policies/{policyId})
	GetPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyId openapi_types.UUID)
	// Put Policy
	// (PUT /v1/{namespaceId}/policies/{policyId})
	PutPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyId openapi_types.UUID)
	// Apply policy
	// (POST /v1/{namespaceId}/policies/{policyId}/apply)
	ApplyPolicyV1(c *gin.Context, namespaceId openapi_types.UUID, policyId openapi_types.UUID)
	// Get namespace profile
	// (GET /v1/{namespaceId}/profile)
	GetNamespaceProfileV1(c *gin.Context, namespaceId openapi_types.UUID)
	// Register namespace
	// (POST /v1/{namespaceId}/profile)
	RegisterNamespaceProfileV1(c *gin.Context, namespaceId openapi_types.UUID)
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

// ListCertificatesV1 operation middleware
func (siw *ServerInterfaceWrapper) ListCertificatesV1(c *gin.Context) {

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

	siw.Handler.ListCertificatesV1(c, namespaceId)
}

// GetCertificateV1 operation middleware
func (siw *ServerInterfaceWrapper) GetCertificateV1(c *gin.Context) {

	var err error

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "namespaceId", c.Param("namespaceId"), &namespaceId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter namespaceId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCertificateV1Params

	headers := c.Request.Header

	// ------------- Optional header parameter "Accept" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Accept")]; found {
		var Accept GetCertificateV1ParamsAccept
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for Accept, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Accept", runtime.ParamLocationHeader, valueList[0], &Accept)
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter Accept: %w", err), http.StatusBadRequest)
			return
		}

		params.Accept = &Accept

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCertificateV1(c, namespaceId, id, params)
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

	siw.Handler.GetPolicyV1(c, namespaceId, policyId)
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

	siw.Handler.PutPolicyV1(c, namespaceId, policyId)
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
	router.GET(options.BaseURL+"/v1/namespaces/:namespaceType", wrapper.ListNamespacesV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/certificates", wrapper.ListCertificatesV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/certificates/:id", wrapper.GetCertificateV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/policies", wrapper.ListPoliciesV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/policies/:policyId", wrapper.GetPolicyV1)
	router.PUT(options.BaseURL+"/v1/:namespaceId/policies/:policyId", wrapper.PutPolicyV1)
	router.POST(options.BaseURL+"/v1/:namespaceId/policies/:policyId/apply", wrapper.ApplyPolicyV1)
	router.GET(options.BaseURL+"/v1/:namespaceId/profile", wrapper.GetNamespaceProfileV1)
	router.POST(options.BaseURL+"/v1/:namespaceId/profile", wrapper.RegisterNamespaceProfileV1)
}
