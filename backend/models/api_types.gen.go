// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oapi-codegen/runtime"
	kmscommon "github.com/stephenzsy/small-kms/backend/common"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for CertificateUsage.
const (
	CertUsageCA         CertificateUsage = "ca"
	CertUsageCARoot     CertificateUsage = "caRoot"
	CertUsageClientAuth CertificateUsage = "clientAuth"
	CertUsageServerAuth CertificateUsage = "serverAuth"
)

// Defines values for CreateProfileRequestType.
const (
	ProfileTypeManagedApplication CreateProfileRequestType = "managed-application"
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

// Defines values for JwtCrv.
const (
	CurveNameP256 JwtCrv = "P-256"
	CurveNameP384 JwtCrv = "P-384"
)

// Defines values for JwtKty.
const (
	KeyTypeEC  JwtKty = "EC"
	KeyTypeRSA JwtKty = "RSA"
)

// Defines values for NamespaceKind.
const (
	NamespaceKindApplication      NamespaceKind = "application"
	NamespaceKindCaInt            NamespaceKind = "ca-int"
	NamespaceKindCaRoot           NamespaceKind = "ca-root"
	NamespaceKindDevice           NamespaceKind = "device"
	NamespaceKindGroup            NamespaceKind = "group"
	NamespaceKindProfile          NamespaceKind = "profile"
	NamespaceKindServicePrincipal NamespaceKind = "service-principal"
	NamespaceKindUser             NamespaceKind = "user"
)

// Defines values for ResourceKind.
const (
	ResourceKindCaInt                 ResourceKind = "ca-int"
	ResourceKindCaRoot                ResourceKind = "ca-root"
	ResourceKindCert                  ResourceKind = "cert"
	ResourceKindCertTemplate          ResourceKind = "cert-template"
	ResourceKindLatestCertForTemplate ResourceKind = "latest-cert-for-template"
	ResourceKindMsGraph               ResourceKind = "ms-graph"
)

// CertificateInfo defines model for CertificateInfo.
type CertificateInfo struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// Id Identifier of the resource
	Id     Identifier      `json:"id"`
	Issuer ResourceLocator `json:"issuer"`

	// Jwk Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	Jwk      JwkProperties          `json:"jwk"`
	Locator  ResourceLocator        `json:"locator"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// NotAfter Expiration date of the certificate
	NotAfter time.Time `json:"notAfter"`

	// NotBefore Expiration date of the certificate
	NotBefore time.Time `json:"notBefore"`
	Pem       *string   `json:"pem,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string          `json:"subjectCommonName"`
	Template          ResourceLocator `json:"template"`

	// Thumbprint X.509 certificate SHA-1 thumbprint
	Thumbprint string `json:"thumbprint"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time         `json:"updated,omitempty"`
	UpdatedBy *string            `json:"updatedBy,omitempty"`
	Usages    []CertificateUsage `json:"usages"`
}

// CertificateInfoFields defines model for CertificateInfoFields.
type CertificateInfoFields struct {
	Issuer ResourceLocator `json:"issuer"`

	// Jwk Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	Jwk JwkProperties `json:"jwk"`

	// NotBefore Expiration date of the certificate
	NotBefore time.Time          `json:"notBefore"`
	Pem       *string            `json:"pem,omitempty"`
	Usages    []CertificateUsage `json:"usages"`
}

// CertificateLifetimeTrigger defines model for CertificateLifetimeTrigger.
type CertificateLifetimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

// CertificateRef defines model for CertificateRef.
type CertificateRef struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// Id Identifier of the resource
	Id       Identifier             `json:"id"`
	Locator  ResourceLocator        `json:"locator"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// NotAfter Expiration date of the certificate
	NotAfter time.Time `json:"notAfter"`

	// SubjectCommonName Common name
	SubjectCommonName string          `json:"subjectCommonName"`
	Template          ResourceLocator `json:"template"`

	// Thumbprint X.509 certificate SHA-1 thumbprint
	Thumbprint string `json:"thumbprint"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// CertificateRefFields defines model for CertificateRefFields.
type CertificateRefFields struct {
	// NotAfter Expiration date of the certificate
	NotAfter time.Time `json:"notAfter"`

	// SubjectCommonName Common name
	SubjectCommonName string          `json:"subjectCommonName"`
	Template          ResourceLocator `json:"template"`

	// Thumbprint X.509 certificate SHA-1 thumbprint
	Thumbprint string `json:"thumbprint"`
}

// CertificateTemplate defines model for CertificateTemplate.
type CertificateTemplate struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// Id Identifier of the resource
	Id             Identifier      `json:"id"`
	IssuerTemplate ResourceLocator `json:"issuerTemplate"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   JwkProperties              `json:"keyProperties"`
	KeyStorePath    *string                    `json:"keyStorePath,omitempty"`
	LifetimeTrigger CertificateLifetimeTrigger `json:"lifetimeTrigger"`
	Locator         ResourceLocator            `json:"locator"`
	Metadata        map[string]interface{}     `json:"metadata,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`

	// Updated Time when the resoruce was last updated
	Updated          *time.Time         `json:"updated,omitempty"`
	UpdatedBy        *string            `json:"updatedBy,omitempty"`
	Usages           []CertificateUsage `json:"usages"`
	ValidityInMonths int32              `json:"validity_months"`
}

// CertificateTemplateFields Certificate fields, may accept template substitutions
type CertificateTemplateFields struct {
	IssuerTemplate ResourceLocator `json:"issuerTemplate"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties    JwkProperties              `json:"keyProperties"`
	KeyStorePath     *string                    `json:"keyStorePath,omitempty"`
	LifetimeTrigger  CertificateLifetimeTrigger `json:"lifetimeTrigger"`
	Usages           []CertificateUsage         `json:"usages"`
	ValidityInMonths int32                      `json:"validity_months"`
}

// CertificateTemplateParameters Certificate fields, may accept template substitutions
type CertificateTemplateParameters struct {
	IssuerTemplate *ResourceLocator `json:"issuerTemplate,omitempty"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   *JwkProperties              `json:"keyProperties,omitempty"`
	KeyStorePath    *string                     `json:"keyStorePath,omitempty"`
	LifetimeTrigger *CertificateLifetimeTrigger `json:"lifetimeTrigger,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string             `json:"subjectCommonName"`
	Usages            []CertificateUsage `json:"usages"`
	ValidityInMonths  *int32             `json:"validity_months,omitempty"`
}

// CertificateTemplateRef defines model for CertificateTemplateRef.
type CertificateTemplateRef struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// Id Identifier of the resource
	Id       Identifier             `json:"id"`
	Locator  ResourceLocator        `json:"locator"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// CertificateTemplateRefFields defines model for CertificateTemplateRefFields.
type CertificateTemplateRefFields struct {
	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`
}

// CertificateUsage defines model for CertificateUsage.
type CertificateUsage string

// CreateManagedApplicationProfileRequest defines model for CreateManagedApplicationProfileRequest.
type CreateManagedApplicationProfileRequest struct {
	Name string                   `json:"name"`
	Type CreateProfileRequestType `json:"type"`
}

// CreateProfileRequest defines model for CreateProfileRequest.
type CreateProfileRequest struct {
	union json.RawMessage
}

// CreateProfileRequestType defines model for CreateProfileRequestType.
type CreateProfileRequestType string

// Identifier Identifier of the resource
type Identifier = kmscommon.Identifier

// IdentifierWithNamespaceKind defines model for IdentifierWithNamespaceKind.
type IdentifierWithNamespaceKind = kmscommon.IdentifierWithKind[NamespaceKind]

// IdentifierWithResourceKind defines model for IdentifierWithResourceKind.
type IdentifierWithResourceKind = kmscommon.IdentifierWithKind[ResourceKind]

// IncludeCertificate defines model for IncludeCertificate.
type IncludeCertificate string

// JwkAlg defines model for JwkAlg.
type JwkAlg string

// KeyOp defines model for JwkKeyOperation.
type KeyOp string

// JwkProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
type JwkProperties struct {
	Alg *JwkAlg `json:"alg,omitempty"`
	Crv *JwtCrv `json:"crv,omitempty"`

	// E RSA exponent
	E     *string `json:"e,omitempty"`
	KeyOp *KeyOp  `json:"key_ops,omitempty"`

	// KeySize RSA key size
	KeySize *int32 `json:"key_size,omitempty"`

	// Kid Key ID
	KeyID *string `json:"kid,omitempty"`
	Kty   JwtKty  `json:"kty"`

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

// JwtCrv defines model for JwtCrv.
type JwtCrv string

// JwtKty defines model for JwtKty.
type JwtKty string

// NamespaceKind defines model for NamespaceKind.
type NamespaceKind string

// Profile defines model for Profile.
type Profile = ProfileRef

// ProfileRef defines model for ProfileRef.
type ProfileRef struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// DisplayName Display name of the resource
	DisplayName string `json:"displayName"`

	// Id Identifier of the resource
	Id Identifier `json:"id"`

	// IsAppManaged Whether the resource is managed by the application
	IsAppManaged *bool                  `json:"isAppManaged,omitempty"`
	Locator      ResourceLocator        `json:"locator"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Type         NamespaceKind          `json:"type"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// ProfileRefFields defines model for ProfileRefFields.
type ProfileRefFields struct {
	// DisplayName Display name of the resource
	DisplayName string `json:"displayName"`

	// IsAppManaged Whether the resource is managed by the application
	IsAppManaged *bool         `json:"isAppManaged,omitempty"`
	Type         NamespaceKind `json:"type"`
}

// ResourceKind defines model for ResourceKind.
type ResourceKind string

// ResourceLocator defines model for ResourceLocator.
type ResourceLocator = kmscommon.Locator[NamespaceKind, ResourceKind]

// ResourceRef defines model for ResourceRef.
type ResourceRef struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// Id Identifier of the resource
	Id       Identifier             `json:"id"`
	Locator  ResourceLocator        `json:"locator"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// CertificateIdPathParameter Identifier of the resource
type CertificateIdPathParameter = Identifier

// CertificateTemplateIdentifierParameter Identifier of the resource
type CertificateTemplateIdentifierParameter = Identifier

// IncludeCertificateParameter defines model for IncludeCertificateParameter.
type IncludeCertificateParameter = IncludeCertificate

// NamespaceIdParameter Identifier of the resource
type NamespaceIdParameter = Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = NamespaceKind

// ProfileIdentifierParameter Identifier of the resource
type ProfileIdentifierParameter = Identifier

// ProfileTypeParameter defines model for ProfileTypeParameter.
type ProfileTypeParameter = NamespaceKind

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = CertificateInfo

// GetCertificateParams defines parameters for GetCertificate.
type GetCertificateParams struct {
	IncludeCertificate    *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
	TemplateId            *Identifier                  `form:"templateId,omitempty" json:"templateId,omitempty"`
	TemplateNamespaceKind *NamespaceKind               `form:"templateNamespaceKind,omitempty" json:"templateNamespaceKind,omitempty"`
	TemplateNamespaceId   *Identifier                  `form:"templateNamespaceId,omitempty" json:"templateNamespaceId,omitempty"`
}

// CreateProfileJSONRequestBody defines body for CreateProfile for application/json ContentType.
type CreateProfileJSONRequestBody = CreateProfileRequest

// PutCertificateTemplateJSONRequestBody defines body for PutCertificateTemplate for application/json ContentType.
type PutCertificateTemplateJSONRequestBody = CertificateTemplateParameters

// AsCreateManagedApplicationProfileRequest returns the union data inside the CreateProfileRequest as a CreateManagedApplicationProfileRequest
func (t CreateProfileRequest) AsCreateManagedApplicationProfileRequest() (CreateManagedApplicationProfileRequest, error) {
	var body CreateManagedApplicationProfileRequest
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCreateManagedApplicationProfileRequest overwrites any union data inside the CreateProfileRequest as the provided CreateManagedApplicationProfileRequest
func (t *CreateProfileRequest) FromCreateManagedApplicationProfileRequest(v CreateManagedApplicationProfileRequest) error {
	v.Type = "managed-application"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCreateManagedApplicationProfileRequest performs a merge with any union data inside the CreateProfileRequest, using the provided CreateManagedApplicationProfileRequest
func (t *CreateProfileRequest) MergeCreateManagedApplicationProfileRequest(v CreateManagedApplicationProfileRequest) error {
	v.Type = "managed-application"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t CreateProfileRequest) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t CreateProfileRequest) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "managed-application":
		return t.AsCreateManagedApplicationProfileRequest()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t CreateProfileRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *CreateProfileRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
