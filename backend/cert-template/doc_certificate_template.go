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
	Digest            string                              `json:"digest"` // checksum of fhte core fields of the template
}

func (d *CertificateTemplateDoc) toModelRef() *models.CertificateTemplateRef {

	return &models.CertificateTemplateRef{
		Id: d.ID.Identifier(),
		Metadata: &models.ResourceMetadata{
			Updated:   utils.ToPtr(d.Updated),
			UpdatedBy: utils.ToPtr(d.UpdatedBy),
			Deleted:   d.Deleted,
		},
		SubjectCommonName: d.SubjectCommonName,
	}
}

func (d *CertificateTemplateDoc) toModel() *models.CertificateTemplate {
	issuerProfileType := models.ProfileTypeIntermediateCA
	if d.IssuerNamespaceID.Kind() == kmsdoc.DocNsTypeCaRoot {
		issuerProfileType = models.ProfileTypeRootCA
	}
	return &models.CertificateTemplate{
		Id: d.ID.Identifier(),
		Metadata: &models.ResourceMetadata{
			Updated:   utils.ToPtr(d.Updated),
			UpdatedBy: utils.ToPtr(d.UpdatedBy),
			Deleted:   d.Deleted,
		},
		SubjectCommonName: d.SubjectCommonName,
		Issuer: &models.CertificateIssuer{
			ProfileType: issuerProfileType,
			ProfileId:   d.IssuerNamespaceID.Identifier(),
			TemplateId:  utils.ToPtr(d.IssuerTemplateID.Identifier()),
		},
		KeyProperties: &models.JwkProperties{
			Alg:     utils.ToPtr(d.KeyProperties.Alg),
			Kty:     d.KeyProperties.Kty,
			Crv:     d.KeyProperties.Crv,
			KeySize: d.KeyProperties.KeySize,
		},
		KeyStorePath:     d.KeyStorePath,
		LifetimeTrigger:  d.LifetimeTrigger,
		Usages:           d.Usages,
		ValidityInMonths: utils.ToPtr(d.ValidityInMonths),
	}
}
