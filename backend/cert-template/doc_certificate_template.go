package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertificateTemplateFlag string

const (
	CertTmplFlagRestrictKtyRsa CertificateTemplateFlag = "kty-rsa"
	CertTmplFlagDelegate       CertificateTemplateFlag = "delegate"
	CertTmplFlagTest           CertificateTemplateFlag = "test"
	CertTmplFlagHasKeyStore    CertificateTemplateFlag = "use-key-store"
	CertTmplFlagKeyExportable  CertificateTemplateFlag = "key-exportable"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	Alg     models.JwkAlg  `json:"alg"`
	Kty     models.JwtKty  `json:"kty"`
	KeySize *int           `json:"key_size,omitempty"`
	Crv     *models.JwtCrv `json:"crv,omitempty"`
	//ReuseKey *bool          `json:"reuse_key,omitempty"`
}

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc

	IssuerNamespaceID kmsdoc.DocNsID                      `json:"issuerNamespaceId"`
	IssuerTemplateID  kmsdoc.DocID                        `json:"issuerTemplateId"`
	Usages            []models.CertificateUsage           `json:"usages"`
	KeyProperties     CertificateTemplateDocKeyProperties `json:"keyProperties"`
	KeyStorePath      *string                             `json:"keyStorePath,omitempty"`
	SubjectCommonName string                              `json:"subjectCn"`
	ValidityInMonths  int32                               `json:"validity_months"`
	LifetimeTrigger   *models.CertificateLifetimeTrigger  `json:"lifetimeTrigger"`
	Digest            []byte                              `json:"version"` // checksum of fhte core fields of the template
}

func (doc *CertificateTemplateDoc) toModel() *models.CertificateTemplate {
	issuerProfileType := models.ProfileTypeIntermediateCA
	if doc.IssuerNamespaceID.Kind() == kmsdoc.DocNsTypeCaRoot {
		issuerProfileType = models.ProfileTypeRootCA
	}
	return &models.CertificateTemplate{
		Issuer: &models.CertificateIssuer{
			ProfileType: issuerProfileType,
			ProfileId:   doc.IssuerNamespaceID.Identifier(),
			TemplateId:  utils.ToPtr(doc.IssuerTemplateID.Identifier()),
		},
		KeyProperties: &models.JwkProperties{
			Alg:     utils.ToPtr(doc.KeyProperties.Alg),
			Kty:     doc.KeyProperties.Kty,
			Crv:     doc.KeyProperties.Crv,
			KeySize: doc.KeyProperties.KeySize,
		},
		KeyStorePath:      doc.KeyStorePath,
		LifetimeTrigger:   doc.LifetimeTrigger,
		SubjectCommonName: doc.SubjectCommonName,
		Usages:            doc.Usages,
		ValidityInMonths:  utils.ToPtr(doc.ValidityInMonths),
	}
}
