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
	externalRef0 "github.com/stephenzsy/small-kms/backend/models"
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

// Defines values for IncludeCertificate.
const (
	IncludeJWK IncludeCertificate = "jwk"
	IncludePEM IncludeCertificate = "pem"
)

// Defines values for NamespaceTypeShortName.
const (
	NSTypeAny              NamespaceTypeShortName = "any"
	NSTypeApplication      NamespaceTypeShortName = "application"
	NSTypeDevice           NamespaceTypeShortName = "device"
	NSTypeGroup            NamespaceTypeShortName = "group"
	NSTypeIntCA            NamespaceTypeShortName = "intermediate-ca"
	NSTypeRootCA           NamespaceTypeShortName = "root-ca"
	NSTypeServicePrincipal NamespaceTypeShortName = "service-principal"
	NSTypeUser             NamespaceTypeShortName = "user"
)

// Defines values for RefType.
const (
	RefTypeCertificate              RefType = "certificate"
	RefTypeCertificateEnrollReceipt RefType = "certificate-enrollment-receipt"
	RefTypeCertificateTemplate      RefType = "certificate-template"
	RefTypeNamespace                RefType = "namespace"
)

// CertificateEnrollmentReceipt defines model for CertificateEnrollmentReceipt.
type CertificateEnrollmentReceipt struct {
	// Expires Time when the enrollment expires
	Expires time.Time `json:"expires"`

	// JwtClaims payload section of the certificate claims, in JWT format, base64url encoded
	JwtClaims string `json:"jwtClaims"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties externalRef0.JwkProperties `json:"keyProperties"`
	Ref           RefWithMetadata            `json:"ref"`

	// RequesterId Unique ID of the user who requested the certificate
	RequesterID openapi_types.UUID `json:"requesterId"`

	// TemplateId Consistent derived ID (UUID v5) of the certificate template
	TemplateID openapi_types.UUID `json:"templateId"`

	// TemplateNamespaceId Unique ID of the namespace of the certificate template
	TemplateNamespaceID openapi_types.UUID `json:"templateNamespaceId"`
}

// CertificateEnrollmentReplyFinalize defines model for CertificateEnrollmentReplyFinalize.
type CertificateEnrollmentReplyFinalize struct {
	// JwtHeader header section of the certificate claims, in JWT format, base64url encoded
	JwtHeader string `json:"jwtHeader"`

	// JwtSignature signature section of the jwt, serves as proof of confirmation finalize the enrollment and being issued a certificate, with header and signature, signed with the key pair with the public key in the same request
	JwtSignature string `json:"jwtSignature"`

	// PublicKey Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	PublicKey externalRef0.JwkProperties `json:"publicKey"`
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
	Jwk *externalRef0.JwkProperties `json:"jwk,omitempty"`

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

// CertificateLifetimeTrigger defines model for CertificateLifetimeTrigger.
type CertificateLifetimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
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
	EmailAddresses []string `json:"emails,omitempty"`
	URIs           []string `json:"uris,omitempty"`
}

// CertificateTemplate defines model for CertificateTemplate.
type CertificateTemplate struct {
	DisplayName string            `json:"displayName"`
	Issuer      CertificateIssuer `json:"issuer"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   *externalRef0.JwkProperties `json:"keyProperties,omitempty"`
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
	KeyProperties   *externalRef0.JwkProperties `json:"keyProperties,omitempty"`
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

// IncludeCertificate defines model for IncludeCertificate.
type IncludeCertificate string

// NamespaceInfo defines model for NamespaceInfo.
type NamespaceInfo struct {
	ObjectType NamespaceTypeShortName `json:"objectType"`
	Ref        RefWithMetadata        `json:"ref"`
}

// NamespaceTypeShortName defines model for NamespaceTypeShortName.
type NamespaceTypeShortName string

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
	Deleted *time.Time `json:"deleted,omitempty"`

	// DisplayName Display name of the object
	DisplayName string             `json:"displayName"`
	ID          openapi_types.UUID `json:"id"`
	IsActive    *bool              `json:"isActive,omitempty"`
	IsDefault   *bool              `json:"isDefault,omitempty"`
	NamespaceID openapi_types.UUID `json:"namespaceId"`
	Type        RefType            `json:"type"`

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

// ServicePrincipalLinkedDevice defines model for ServicePrincipalLinkedDevice.
type ServicePrincipalLinkedDevice struct {
	ApplicationClientID openapi_types.UUID `json:"applicationClientId"`

	// ApplicationOid Object ID of the application
	ApplicationOID openapi_types.UUID `json:"applicationOid"`
	DeviceID       openapi_types.UUID `json:"deviceId"`

	// DeviceOid Object ID of the device
	DeviceOID          openapi_types.UUID `json:"deviceOid"`
	Ref                RefWithMetadata    `json:"ref"`
	ServicePrincipalID openapi_types.UUID `json:"servicePrincipalId"`
	Status             string             `json:"status"`
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
type ErrorResponse map[string]interface{}

// RefListResponse defines model for RefListResponse.
type RefListResponse = []RefWithMetadata

// ListNamespacesByTypeV2Params defines parameters for ListNamespacesByTypeV2.
type ListNamespacesByTypeV2Params struct {
	NamespaceType NamespaceTypeParameter `form:"namespaceType" json:"namespaceType"`
}

// ListCertificateTemplatesV2Params defines parameters for ListCertificateTemplatesV2.
type ListCertificateTemplatesV2Params struct {
	IncludeDefaultForType *NamespaceTypeShortName `form:"includeDefaultForType,omitempty" json:"includeDefaultForType,omitempty"`
}

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

// PutCertificateTemplateV2JSONRequestBody defines body for PutCertificateTemplateV2 for application/json ContentType.
type PutCertificateTemplateV2JSONRequestBody = CertificateTemplateParameters

// BeginEnrollCertificateV2JSONRequestBody defines body for BeginEnrollCertificateV2 for application/json ContentType.
type BeginEnrollCertificateV2JSONRequestBody = CertificateEnrollmentRequest

// CompleteCertificateEnrollmentV2JSONRequestBody defines body for CompleteCertificateEnrollmentV2 for application/json ContentType.
type CompleteCertificateEnrollmentV2JSONRequestBody = CertificateEnrollmentReplyFinalize

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
