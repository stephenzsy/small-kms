package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateTemplateFlag string

const (
	CertTmplFlagRestrictKtyRsa CertificateTemplateFlag = "kty-rsa"
	CertTmplFlagDelegate       CertificateTemplateFlag = "delegate"
	CertTmplFlagTest           CertificateTemplateFlag = "test"
	CertTmplFlagHasKeyStore    CertificateTemplateFlag = "use-key-store"
	CertTmplFlagKeyExportable  CertificateTemplateFlag = "key-exportable"
)

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc

	IssuerTemplate    models.ResourceLocator            `json:"issuerTemplate"`
	Usages            []models.CertificateUsage         `json:"usages"`
	KeySpec           CertKeySpec                       `json:"keySpec"`
	KeyStorePath      *string                           `json:"keyStorePath,omitempty"`
	SubjectCommonName string                            `json:"subjectCn"`
	ValidityInMonths  int32                             `json:"validity_months"`
	LifetimeTrigger   models.CertificateLifetimeTrigger `json:"lifetimeTrigger"`
	Digest            kmsdoc.HexStringStroable          `json:"digest"` // checksum of fhte core fields of the template
}

func (d *CertificateTemplateDoc) populateRef(r *models.CertificateTemplateRefComposed) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.SubjectCommonName = d.SubjectCommonName
}

func (d *CertificateTemplateDoc) toModelRef() (r *models.CertificateTemplateRefComposed) {
	r = new(models.CertificateTemplateRefComposed)
	d.populateRef(r)
	return
}

func (d *CertificateTemplateDoc) toModel() *models.CertificateTemplateComposed {
	r := new(models.CertificateTemplateComposed)
	d.populateRef(&r.CertificateTemplateRefComposed)
	r.IssuerTemplate = d.IssuerTemplate
	d.KeySpec.PopulateKeyProperties(&r.KeyProperties)
	r.KeyStorePath = d.KeyStorePath
	r.LifetimeTrigger = d.LifetimeTrigger
	r.ValidityInMonths = d.ValidityInMonths
	r.Usages = d.Usages
	return r
}
