// Package admin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for CertificateEnrollmentTargetType.
const (
	CertEnrollTargetTypeDeviceLinkedServicePrincipal CertificateEnrollmentTargetType = "device-linked-service-principal"
)

// Defines values for CertificateUsage.
const (
	UsageAADClientCredential CertificateUsage = "aad-client-credential"
	UsageClientOnly          CertificateUsage = "client-only"
	UsageIntCA               CertificateUsage = "intermediate-ca"
	UsageRootCA              CertificateUsage = "root-ca"
	UsageServerAndClient     CertificateUsage = "server-and-client"
	UsageServerOnly          CertificateUsage = "server-only"
)

// Defines values for CurveName.
const (
	CurveNameP256 CurveName = "P-256"
	CurveNameP384 CurveName = "P-384"
)

// Defines values for IncludeCertificate.
const (
	IncludeJWK IncludeCertificate = "jwk"
	IncludePEM IncludeCertificate = "pem"
)

// Defines values for JwkAlg.
const (
	AlgES256 JwkAlg = "ES256"
	AlgES384 JwkAlg = "ES384"
	AlgRS256 JwkAlg = "RS256"
	AlgRS384 JwkAlg = "RS384"
	AlgRS512 JwkAlg = "RS512"
)

// Defines values for KeyOp.
const (
	KeyOpDecrypt   KeyOp = "decrypt"
	KeyOpEncrypt   KeyOp = "encrypt"
	KeyOpSign      KeyOp = "sign"
	KeyOpUnwrapKey KeyOp = "unwrapKey"
	KeyOpVerify    KeyOp = "verify"
	KeyOpWrapKey   KeyOp = "wrapKey"
)

// Defines values for KeySize.
const (
	KeySize2048 KeySize = 2048
	KeySize3072 KeySize = 3072
	KeySize4096 KeySize = 4096
)

// Defines values for KeyType.
const (
	KeyTypeEC  KeyType = "EC"
	KeyTypeRSA KeyType = "RSA"
)

// Defines values for NamespaceType.
const (
	NamespaceTypeBuiltInCaInt            NamespaceType = "#builtin.ca.intermediate"
	NamespaceTypeBuiltInCaRoot           NamespaceType = "#builtin.ca.root"
	NamespaceTypeMsGraphApplication      NamespaceType = "#microsoft.graph.application"
	NamespaceTypeMsGraphDevice           NamespaceType = "#microsoft.graph.device"
	NamespaceTypeMsGraphGroup            NamespaceType = "#microsoft.graph.group"
	NamespaceTypeMsGraphServicePrincipal NamespaceType = "#microsoft.graph.servicePrincipal"
	NamespaceTypeMsGraphUser             NamespaceType = "#microsoft.graph.user"
)

// Defines values for NamespaceTypeShortName.
const (
	NSTypeApplication      NamespaceTypeShortName = "application"
	NSTypeDevice           NamespaceTypeShortName = "device"
	NSTypeGroup            NamespaceTypeShortName = "group"
	NSTypeIntCA            NamespaceTypeShortName = "intermediate-ca"
	NSTypeRootCA           NamespaceTypeShortName = "root-ca"
	NSTypeServicePrincipal NamespaceTypeShortName = "service-principal"
	NSTypeUnknown          NamespaceTypeShortName = "unknown"
	NSTypeUser             NamespaceTypeShortName = "user"
)

// Defines values for PolicyType.
const (
	PolicyTypeCertEnroll PolicyType = "certEnroll"
)

// Defines values for RefType.
const (
	RefTypeCertificate         RefType = "certificate"
	RefTypeCertificateTemplate RefType = "certificate-template"
	RefTypeNamespace           RefType = "namespace"
)

// ApplyPolicyRequest defines model for ApplyPolicyRequest.
type ApplyPolicyRequest struct {
	// CheckConsistency Check consistency of the policy
	CheckConsistency *bool `json:"checkConsistency,omitempty"`

	// ForceRenewCertificate Force certificate renewal
	ForceRenewCertificate *bool `json:"forceRenewCertificate,omitempty"`
}

// CertificateEnrollPolicyParameters defines model for CertificateEnrollPolicyParameters.
type CertificateEnrollPolicyParameters struct {
	AllowedUsages       []CertificateUsage `json:"allowedUsages"`
	MaxValidityInMonths int32              `json:"maxValidityInMonths"`
}

// CertificateEnrollmentReceipt defines model for CertificateEnrollmentReceipt.
type CertificateEnrollmentReceipt struct {
	// JwtPayload payload section of the certificate claims, in JWT format, base64url encoded
	JwtPayload *string         `json:"jwtPayload,omitempty"`
	Ref        RefWithMetadata `json:"ref"`

	// TemplateId Consistent derived ID (UUID v5) of the certificate template
	TemplateID *openapi_types.UUID `json:"templateId,omitempty"`
}

// CertificateEnrollmentRequest defines model for CertificateEnrollmentRequest.
type CertificateEnrollmentRequest struct {
	union json.RawMessage
}

// CertificateEnrollmentRequestDeviceLinkedServicePrincipal defines model for CertificateEnrollmentRequestDeviceLinkedServicePrincipal.
type CertificateEnrollmentRequestDeviceLinkedServicePrincipal struct {
	// AppId Client ID of the application
	AppID openapi_types.UUID `json:"appId"`

	// CommonName Common Name to appear in the certificate
	CommonName *string `json:"commonName,omitempty"`

	// DeviceNamespaceId Object ID of the device
	DeviceNamespaceID openapi_types.UUID `json:"deviceNamespaceId"`

	// LinkId Unique ID of the device link
	DeviceLinkID openapi_types.UUID `json:"linkId"`

	// ServicePrincipalId Object ID of the service principal
	ServicePrincipalID openapi_types.UUID              `json:"servicePrincipalId"`
	Type               CertificateEnrollmentTargetType `json:"type"`
}

// CertificateEnrollmentTargetType defines model for CertificateEnrollmentTargetType.
type CertificateEnrollmentTargetType string

// CertificateInfo defines model for CertificateInfo.
type CertificateInfo struct {
	// CommonName Common name
	CommonName        string `json:"commonName"`
	IssuerCertificate Ref    `json:"issuerCertificate"`

	// Jwk Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	Jwk *JwkProperties `json:"jwk,omitempty"`

	// NotAfter Expiration date of the certificate
	NotAfter time.Time `json:"notAfter"`

	// NotBefore Expiration date of the certificate
	NotBefore               time.Time                           `json:"notBefore"`
	Pem                     *string                             `json:"pem,omitempty"`
	Ref                     RefWithMetadata                     `json:"ref"`
	Subject                 string                              `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Template                Ref                                 `json:"template"`
	Usage                   CertificateUsage                    `json:"usage"`
}

// CertificateIssuer defines model for CertificateIssuer.
type CertificateIssuer struct {
	NamespaceID   openapi_types.UUID     `json:"namespaceId"`
	NamespaceType NamespaceTypeShortName `json:"namespaceType"`

	// TemplateId if certificate ID is not specified, use template ID to find the latest certificate, use default value if not specified
	TemplateID *openapi_types.UUID `json:"templateId,omitempty"`
}

// CertificateIssuerParameters defines model for CertificateIssuerParameters.
type CertificateIssuerParameters struct {
	// IssuerNamespaceId ID of the issuer namespace
	IssuerNamespaceID openapi_types.UUID `json:"issuerNamespaceId"`

	// IssuerPolicyIdentifier ID of the issuer policy
	IssuerPolicyIdentifier *string `json:"issuerPolicyIdentifier,omitempty"`
}

// CertificateLifetimeTrigger defines model for CertificateLifetimeTrigger.
type CertificateLifetimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

// CertificateRequestPolicyParameters defines model for CertificateRequestPolicyParameters.
type CertificateRequestPolicyParameters struct {
	// IssuerNamespaceId ID of the issuer namespace
	IssuerNamespaceID openapi_types.UUID `json:"issuerNamespaceId"`

	// IssuerPolicyIdentifier ID of the issuer policy
	IssuerPolicyIdentifier  *string                             `json:"issuerPolicyIdentifier,omitempty"`
	KeyProperties           *KeyProperties                      `json:"keyProperties,omitempty"`
	LifetimeTrigger         *CertificateLifetimeTrigger         `json:"lifetimeTrigger,omitempty"`
	Subject                 CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                    `json:"usage"`
	ValidityInMonths        *int32                              `json:"validity_months,omitempty"`
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
	DNSNames       []string `json:"dns_names,omitempty"`
	EmailAddresses []string `json:"emails,omitempty"`
	IPAddresses    []string `json:"ipAddrs,omitempty"`
	URIs           []string `json:"uris,omitempty"`
}

// CertificateTemplate defines model for CertificateTemplate.
type CertificateTemplate struct {
	DisplayName string            `json:"displayName"`
	Issuer      CertificateIssuer `json:"issuer"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   *JwkProperties              `json:"keyProperties,omitempty"`
	KeyStorePath    *string                     `json:"keyStorePath,omitempty"`
	LifetimeTrigger *CertificateLifetimeTrigger `json:"lifetimeTrigger,omitempty"`
	Ref             RefWithMetadata             `json:"ref"`

	// ReuseKey Keep using the same key version if exists
	ReuseKey                *bool                               `json:"reuse_key,omitempty"`
	Subject                 CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                    `json:"usage"`
	ValidityInMonths        *int32                              `json:"validity_months,omitempty"`
}

// CertificateTemplateParameters Certificate fields, may accept template substitutions
type CertificateTemplateParameters struct {
	DisplayName string            `json:"displayName"`
	Issuer      CertificateIssuer `json:"issuer"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   *JwkProperties              `json:"keyProperties,omitempty"`
	KeyStorePath    *string                     `json:"keyStorePath,omitempty"`
	LifetimeTrigger *CertificateLifetimeTrigger `json:"lifetimeTrigger,omitempty"`

	// ReuseKey Keep using the same key version if exists
	ReuseKey                *bool                               `json:"reuse_key,omitempty"`
	Subject                 CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                    `json:"usage"`
	ValidityInMonths        *int32                              `json:"validity_months,omitempty"`
}

// CertificateUsage defines model for CertificateUsage.
type CertificateUsage string

// CurveName defines model for CurveName.
type CurveName string

// IncludeCertificate defines model for IncludeCertificate.
type IncludeCertificate string

// JwkAlg defines model for JwkAlg.
type JwkAlg string

// KeyOp defines model for JwkKeyOperation.
type KeyOp string

// KeySize defines model for JwkKeySize.
type KeySize int32

// JwkProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
type JwkProperties struct {
	Alg *JwkAlg    `json:"alg,omitempty"`
	Crv *CurveName `json:"crv,omitempty"`

	// E RSA exponent
	E       *string  `json:"e,omitempty"`
	KeyOp   *KeyOp   `json:"key_ops,omitempty"`
	KeySize *KeySize `json:"key_size,omitempty"`

	// Kid Key ID
	KeyID *string `json:"kid,omitempty"`
	Kty   KeyType `json:"kty"`

	// N RSA modulus
	N *string `json:"n,omitempty"`

	// X EC x coordinate
	X *string `json:"x,omitempty"`

	// X5c X.509 certificate chain
	CertificateChain []string `json:"x5c,omitempty"`

	// X5t X.509 certificate SHA-1 thumbprint
	CertificateThumbprint *string `json:"x5t,omitempty"`

	// X5tS256 X.509 certificate SHA-256 thumbprint
	CertificateThumbprintSHA256 *string `json:"x5t#S256,omitempty"`

	// X5u X.509 certificate URL
	CertificateURL *string `json:"x5u,omitempty"`

	// Y EC y coordinate
	Y *string `json:"y,omitempty"`
}

// KeyProperties defines model for KeyProperties.
type KeyProperties struct {
	Crv     *CurveName `json:"crv,omitempty"`
	KeySize *KeySize   `json:"key_size,omitempty"`
	Kty     KeyType    `json:"kty"`

	// ReuseKey Keep using the same key version if exists
	ReuseKey *bool `json:"reuse_key,omitempty"`
}

// KeyType defines model for KeyType.
type KeyType string

// NamespaceInfo defines model for NamespaceInfo.
type NamespaceInfo struct {
	ObjectType NamespaceType   `json:"objectType"`
	Ref        RefWithMetadata `json:"ref"`
}

// NamespaceProfile defines model for NamespaceProfile.
type NamespaceProfile struct {
	// AppId \#microsoft.graph.application appId
	AppID *string `json:"appId,omitempty"`

	// Deleted Time when the policy was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// DeviceId \#microsoft.graph.device deviceId
	DeviceID        *string            `json:"deviceId,omitempty"`
	DeviceOwnership *string            `json:"deviceOwnership,omitempty"`
	DisplayName     string             `json:"displayName"`
	ID              openapi_types.UUID `json:"id"`

	// IsCompliant \#microsoft.graph.device isCompliant
	IsCompliant *bool `json:"isCompliant,omitempty"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	ObjectType  NamespaceType      `json:"objectType"`

	// OperatingSystem \#microsoft.graph.device operatingSystem
	OperatingSystem *string `json:"operatingSystem,omitempty"`

	// OperatingSystemVersion \#microsoft.graph.device operatingSystemVersion
	OperatingSystemVersion *string `json:"operatingSystemVersion,omitempty"`
	ServicePrincipalType   *string `json:"servicePrincipalType,omitempty"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy         string  `json:"updatedBy"`
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
}

// NamespaceRef defines model for NamespaceRef.
type NamespaceRef struct {
	// Deleted Time when the policy was deleted
	Deleted     *time.Time         `json:"deleted,omitempty"`
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

// NamespaceTypeShortName defines model for NamespaceTypeShortName.
type NamespaceTypeShortName string

// Policy defines model for Policy.
type Policy struct {
	CertEnroll  *CertificateEnrollPolicyParameters  `json:"certEnroll,omitempty"`
	CertRequest *CertificateRequestPolicyParameters `json:"certRequest,omitempty"`

	// Deleted Time when the policy was deleted
	Deleted *time.Time         `json:"deleted,omitempty"`
	ID      openapi_types.UUID `json:"id"`

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
	CertEnroll  *CertificateEnrollPolicyParameters  `json:"certEnroll,omitempty"`
	CertRequest *CertificateRequestPolicyParameters `json:"certRequest,omitempty"`
	PolicyType  PolicyType                          `json:"policyType"`
}

// PolicyRef defines model for PolicyRef.
type PolicyRef struct {
	// Deleted Time when the policy was deleted
	Deleted *time.Time         `json:"deleted,omitempty"`
	ID      openapi_types.UUID `json:"id"`

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

	// Deleted Time when the policy was deleted
	Deleted *time.Time         `json:"deleted,omitempty"`
	ID      openapi_types.UUID `json:"id"`
	Message string             `json:"message"`

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

// Ref defines model for Ref.
type Ref struct {
	ID          openapi_types.UUID `json:"id"`
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	Type        RefType            `json:"type"`
}

// RefType defines model for RefType.
type RefType string

// RefWithMetadata defines model for RefWithMetadata.
type RefWithMetadata struct {
	// Deleted Time when the object was deleted
	Deleted       *time.Time             `json:"deleted,omitempty"`
	ID            openapi_types.UUID     `json:"id"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	NamespaceID   openapi_types.UUID     `json:"namespaceId"`
	NamespaceType NamespaceTypeShortName `json:"namespaceType"`
	Type          RefType                `json:"type"`

	// Updated Time when the object was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who last updated the object
	UpdatedBy string `json:"updatedBy"`
}

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
	// Deleted Time when the policy was deleted
	Deleted *time.Time         `json:"deleted,omitempty"`
	ID      openapi_types.UUID `json:"id"`

	// NamespaceId Unique ID of the namespace
	NamespaceID openapi_types.UUID `json:"namespaceId"`

	// Updated Time when the policy was last updated
	Updated time.Time `json:"updated"`

	// UpdatedBy Unique ID of the user who created the policy
	UpdatedBy string `json:"updatedBy"`
}

// ServicePrincipalLinkedDevice defines model for ServicePrincipalLinkedDevice.
type ServicePrincipalLinkedDevice struct {
	ApplicationClientID openapi_types.UUID `json:"applicationClientId"`

	// ApplicationOid Object ID of the application
	ApplicationOID openapi_types.UUID `json:"applicationOid"`
	DeviceID       openapi_types.UUID `json:"deviceId"`

	// DeviceOid Object ID of the device
	DeviceOID          openapi_types.UUID `json:"deviceOid"`
	ServicePrincipalID openapi_types.UUID `json:"servicePrincipalId"`
}

// CertIdParameter defines model for CertIdParameter.
type CertIdParameter = openapi_types.UUID

// IncludeCertificateParameter defines model for IncludeCertificateParameter.
type IncludeCertificateParameter = IncludeCertificate

// NamespaceIdParameter defines model for NamespaceIdParameter.
type NamespaceIdParameter = openapi_types.UUID

// NamespaceTypeParameter defines model for NamespaceTypeParameter.
type NamespaceTypeParameter = NamespaceTypeShortName

// TemplateIdParameter defines model for TemplateIdParameter.
type TemplateIdParameter = openapi_types.UUID

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = CertificateInfo

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Code                 *string                `json:"code,omitempty"`
	Message              *string                `json:"message,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// RefListResponse defines model for RefListResponse.
type RefListResponse = []RefWithMetadata

// IssueCertificateByTemplateV2Params defines parameters for IssueCertificateByTemplateV2.
type IssueCertificateByTemplateV2Params struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
}

// GetLatestCertificateByTemplateV2Params defines parameters for GetLatestCertificateByTemplateV2.
type GetLatestCertificateByTemplateV2Params struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
}

// GetCertificateV2Params defines parameters for GetCertificateV2.
type GetCertificateV2Params struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
}

// CompleteCertificateEnrollmentV2Params defines parameters for CompleteCertificateEnrollmentV2.
type CompleteCertificateEnrollmentV2Params struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
}

// ApplyPolicyV1JSONRequestBody defines body for ApplyPolicyV1 for application/json ContentType.
type ApplyPolicyV1JSONRequestBody = ApplyPolicyRequest

// BeginEnrollCertificateV2JSONRequestBody defines body for BeginEnrollCertificateV2 for application/json ContentType.
type BeginEnrollCertificateV2JSONRequestBody = CertificateEnrollmentRequest

// PutCertificateTemplateV2JSONRequestBody defines body for PutCertificateTemplateV2 for application/json ContentType.
type PutCertificateTemplateV2JSONRequestBody = CertificateTemplateParameters

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

// AsCertificateEnrollmentRequestDeviceLinkedServicePrincipal returns the union data inside the CertificateEnrollmentRequest as a CertificateEnrollmentRequestDeviceLinkedServicePrincipal
func (t CertificateEnrollmentRequest) AsCertificateEnrollmentRequestDeviceLinkedServicePrincipal() (CertificateEnrollmentRequestDeviceLinkedServicePrincipal, error) {
	var body CertificateEnrollmentRequestDeviceLinkedServicePrincipal
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCertificateEnrollmentRequestDeviceLinkedServicePrincipal overwrites any union data inside the CertificateEnrollmentRequest as the provided CertificateEnrollmentRequestDeviceLinkedServicePrincipal
func (t *CertificateEnrollmentRequest) FromCertificateEnrollmentRequestDeviceLinkedServicePrincipal(v CertificateEnrollmentRequestDeviceLinkedServicePrincipal) error {
	v.Type = "device-linked-service-principal"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCertificateEnrollmentRequestDeviceLinkedServicePrincipal performs a merge with any union data inside the CertificateEnrollmentRequest, using the provided CertificateEnrollmentRequestDeviceLinkedServicePrincipal
func (t *CertificateEnrollmentRequest) MergeCertificateEnrollmentRequestDeviceLinkedServicePrincipal(v CertificateEnrollmentRequestDeviceLinkedServicePrincipal) error {
	v.Type = "device-linked-service-principal"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t CertificateEnrollmentRequest) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t CertificateEnrollmentRequest) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "device-linked-service-principal":
		return t.AsCertificateEnrollmentRequestDeviceLinkedServicePrincipal()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t CertificateEnrollmentRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *CertificateEnrollmentRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
