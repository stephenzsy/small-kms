package admin

import (
	"crypto/x509/pkix"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	// Kty      KeyType       `json:"kty"`
	// KeySize  *KeySize      `json:"key_size,omitempty"`
	// Crv      *CurveName    `json:"crv,omitempty"`
	ReuseKey *bool `json:"reuse_key,omitempty"`
}

type CertificateTemplateDocLifeTimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

type CertificateTemplateDocSubject struct {
	// CertificateSubject
	cachedString *string
}

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc
	DisplayName       string    `json:"displayName"`
	IssuerNamespaceID uuid.UUID `json:"issuerNamespaceId"`
	// IssuerNameSpaceType     NamespaceTypeShortName                `json:"issuerNameSpaceType"`
	IssuerTemplateID        kmsdoc.KmsDocID                       `json:"issuerTemplateId"`
	KeyProperties           CertificateTemplateDocKeyProperties   `json:"keyProperties"`
	KeyStorePath            *string                               `json:"keyStorePath,omitempty"`
	Subject                 CertificateTemplateDocSubject         `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames   `json:"sans,omitempty"`
	Usage                   CertificateUsage                      `json:"usage"`
	ValidityInMonths        int32                                 `json:"validity_months"`
	LifetimeTrigger         CertificateTemplateDocLifeTimeTrigger `json:"lifetimeTrigger"`
	VariablesEnabeled       bool                                  `json:"variablesEnabeled"`
}

func (doc *CertificateTemplateDoc) IsActive() bool {
	return doc.Deleted == nil || doc.Deleted.IsZero()
}

func (doc *CertificateTemplateDoc) IssuerCertificateDocID() kmsdoc.KmsDocID {
	return kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForTemplate, doc.IssuerTemplateID.GetUUID())
}

func (s *CertificateTemplateDocSubject) pkixName() (name pkix.Name) {
	// name.CommonName = s.CN
	// if s.C != nil && len(*s.C) > 0 {
	// 	name.Country = []string{*s.C}
	// }
	// if s.O != nil && len(*s.O) > 0 {
	// 	name.Organization = []string{*s.O}
	// }
	// if s.OU != nil && len(*s.OU) > 0 {
	// 	name.OrganizationalUnit = []string{*s.OU}
	// }
	return
}

func (s *CertificateTemplateDocSubject) String() string {
	if s == nil {
		return ""
	}
	if s.cachedString != nil {
		return *s.cachedString
	}
	name := s.pkixName()
	str := name.String()
	s.cachedString = &str
	return str
}
