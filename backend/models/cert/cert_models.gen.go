// Package certmodels provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package certmodels

import (
	"net"

	externalRef0 "github.com/stephenzsy/small-kms/backend/models"
	externalRef1 "github.com/stephenzsy/small-kms/backend/models/key"
)

// Defines values for CertificateFlag.
const (
	CertificateFlagClientAuth CertificateFlag = "clientAuth"
	CertificateFlagServerAuth CertificateFlag = "serverAuth"
)

// Defines values for CertificateStatus.
const (
	CertificateStatusDeactivated     CertificateStatus = "deactivated"
	CertificateStatusIssued          CertificateStatus = "issued"
	CertificateStatusPending         CertificateStatus = "pending"
	CertificateStatusPendingExternal CertificateStatus = "pending-external"
	CertificateStatusUnverified      CertificateStatus = "unverified"
)

// Certificate defines model for Certificate.
type Certificate = certificateComposed

// CertificateExternalIssuer defines model for CertificateExternalIssuer.
type CertificateExternalIssuer = certificateExternalIssuerComposed

// CertificateExternalIssuerAcme defines model for CertificateExternalIssuerAcme.
type CertificateExternalIssuerAcme struct {
	AccountKeyID           string   `json:"accountKeyId"`
	AccountURL             string   `json:"accountUrl"`
	AzureDNSZoneResourceID string   `json:"azureDnsZoneResourceId"`
	Contacts               []string `json:"contacts"`
	DirectoryURL           string   `json:"directoryUrl"`
}

// CertificateExternalIssuerFields defines model for CertificateExternalIssuerFields.
type CertificateExternalIssuerFields struct {
	Acme *CertificateExternalIssuerAcme `json:"acme,omitempty"`
}

// CertificateFields defines model for CertificateFields.
type CertificateFields struct {
	// Cid Key Vault certificate ID
	KeyVaultCertificateID string                   `json:"cid,omitempty"`
	Flags                 []CertificateFlag        `json:"flags,omitempty"`
	Identififier          string                   `json:"identififier"`
	IssuerIdentifier      string                   `json:"issuerIdentifier"`
	Jwk                   *externalRef1.JsonWebKey `json:"jwk,omitempty"`
	Nbf                   externalRef0.NumericDate `json:"nbf"`
	OneTimePkcs12Key      *externalRef1.JsonWebKey `json:"oneTimePkcs12Key,omitempty"`

	// Sid Key Vault Secret ID
	KeyVaultSecretID        string                   `json:"sid,omitempty"`
	Subject                 string                   `json:"subject"`
	SubjectAlternativeNames *SubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
}

// CertificateFlag defines model for CertificateFlag.
type CertificateFlag string

// CertificatePolicy defines model for CertificatePolicy.
type CertificatePolicy = certificatePolicyComposed

// CertificatePolicyFields defines model for CertificatePolicyFields.
type CertificatePolicyFields struct {
	AllowEnroll   bool              `json:"allowEnroll"`
	AllowGenerate bool              `json:"allowGenerate"`
	ExpiryTime    string            `json:"expiryTime"`
	Flags         []CertificateFlag `json:"flags,omitempty"`

	// IssuerPolicyIdentifier Policy identififier of parent issuer
	IssuerPolicyIdentifier string `json:"issuerPolicyIdentifier"`
	KeyExportable          bool   `json:"keyExportable"`

	// KeySpec these attributes should mostly confirm to JWK (RFC7517)
	KeySpec                 externalRef1.JsonWebKeySpec `json:"keySpec"`
	Subject                 CertificateSubject          `json:"subject"`
	SubjectAlternativeNames *SubjectAlternativeNames    `json:"subjectAlternativeNames,omitempty"`
}

// CertificatePolicyParameters defines model for CertificatePolicyParameters.
type CertificatePolicyParameters struct {
	AllowEnroll            *bool             `json:"allowEnroll,omitempty"`
	AllowGenerate          *bool             `json:"allowGenerate,omitempty"`
	DisplayName            string            `json:"displayName,omitempty"`
	ExpiryTime             string            `json:"expiryTime,omitempty"`
	Flags                  []CertificateFlag `json:"flags,omitempty"`
	IssuerPolicyIdentifier string            `json:"issuerPolicyIdentifier,omitempty"`
	KeyExportable          *bool             `json:"keyExportable,omitempty"`

	// KeySpec these attributes should mostly confirm to JWK (RFC7517)
	KeySpec                 *externalRef1.JsonWebKeySpec `json:"keySpec,omitempty"`
	Subject                 CertificateSubject           `json:"subject"`
	SubjectAlternativeNames *SubjectAlternativeNames     `json:"subjectAlternativeNames,omitempty"`
}

// CertificateRef defines model for CertificateRef.
type CertificateRef = certificateRefComposed

// CertificateRefFields defines model for CertificateRefFields.
type CertificateRefFields struct {
	Exp              externalRef0.NumericDate  `json:"exp"`
	Iat              *externalRef0.NumericDate `json:"iat,omitempty"`
	PolicyIdentifier string                    `json:"policyIdentifier"`
	Status           CertificateStatus         `json:"status"`

	// Thumbprint Hex encoded certificate thumbprint
	Thumbprint string `json:"thumbprint"`
}

// CertificateStatus defines model for CertificateStatus.
type CertificateStatus string

// CertificateSubject defines model for CertificateSubject.
type CertificateSubject struct {
	CommonName string `json:"cn"`
}

// EnrollCertificateRequest defines model for EnrollCertificateRequest.
type EnrollCertificateRequest struct {
	PublicKey            externalRef1.JsonWebKey `json:"publicKey"`
	WithOneTimePkcs12Key *bool                   `json:"withOneTimePkcs12Key,omitempty"`
}

// ExchangePKCS12Request defines model for ExchangePKCS12Request.
type ExchangePKCS12Request struct {
	// Legacy Use legacy PKCS12 cipher
	Legacy *bool `json:"legacy,omitempty"`

	// PasswordProtected Encrypt the PKCS12 file with a generated password
	PasswordProtected bool `json:"passwordProtected"`

	// Payload JWE encrypted private key in JWK
	Payload string `json:"payload"`
}

// ExchangePKCS12Result defines model for ExchangePKCS12Result.
type ExchangePKCS12Result struct {
	// Password Password used to encrypt the PKCS12 file
	Password string `json:"password"`

	// Payload JWE encrypted PKCS12 file, encrypted with the symmetric key from the request
	Payload string `json:"payload"`
}

// SubjectAlternativeNames defines model for SubjectAlternativeNames.
type SubjectAlternativeNames struct {
	DNSNames    []string `json:"dnsNames,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	IPAddresses []net.IP `json:"ipAddresses,omitempty"`
}

// CertificateExternalIssuerResponse defines model for CertificateExternalIssuerResponse.
type CertificateExternalIssuerResponse = CertificateExternalIssuer

// CertificatePolicyResponse defines model for CertificatePolicyResponse.
type CertificatePolicyResponse = CertificatePolicy

// CertificateRefsResponse defines model for CertificateRefsResponse.
type CertificateRefsResponse = []CertificateRef

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = Certificate
