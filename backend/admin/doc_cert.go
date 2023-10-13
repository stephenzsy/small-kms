package admin

import (
	"time"

	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertDoc struct {
	kmsdoc.BaseDoc

	// alias for certs with L prefix
	AliasID *kmsdoc.KmsDocID `json:"aliasId,omitempty"`

	IssuerNamespaceID   uuid.UUID       `json:"issuerNamespaceId"`
	IssuerCertificateID kmsdoc.KmsDocID `json:"issuerCertId"`
	TemplateID          kmsdoc.KmsDocID `json:"templateId"`
	Subject             string          `json:"subject"`
	SubjectBase         string          `json:"subjectBase"`
	// KeyInfo                 JwkProperties                       `json:"keyInfo"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"sans,omitempty"`
	NotBefore               time.Time                           `json:"notBefore"`
	NotAfter                time.Time                           `json:"notAfter"`
	CertStorePath           string                              `json:"certStorePath"` // certificate storage path in blob storage
	CommonName              string                              `json:"name"`
	Usage                   CertificateUsage                    `json:"usage"`
	FingerprintSHA1Hex      string                              `json:"fingerprint"` // information only
}

func (doc *CertDoc) IsActive() bool {
	if doc == nil {
		return false
	}
	if doc.Deleted != nil && !doc.Deleted.IsZero() {
		return false
	}
	now := time.Now()
	if now.After(doc.NotAfter) || now.Before(doc.NotBefore) {
		return false
	}
	return true
}
