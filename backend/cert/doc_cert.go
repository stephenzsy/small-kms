package cert

import (
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateStatus string

const (
	CertStatusInitial CertificateStatus = ""
	CertStatusPending CertificateStatus = "pending"
	CertStatusIssued  CertificateStatus = "issued"
)

type ResourceLocator = models.ResourceLocator

type CertDoc struct {
	kmsdoc.BaseDoc

	Status CertificateStatus `json:"status"` // certificate status

	// X509 certificate info
	SerialNumber      SerialNumberStorable      `json:"serialNumber"`
	SubjectCommonName string                    `json:"subjectCommonName"`
	NotBefore         kmsdoc.TimeStorable       `json:"notBefore"`
	NotAfter          kmsdoc.TimeStorable       `json:"notAfter"`
	Usages            []models.CertificateUsage `json:"usages"`
	KeySpec           ct.CertKeySpec            `json:"keySpec"`
	KeyStorePath      *string                   `json:"keyStorePath,omitempty"`
	CertStorePath     string                    `json:"certStorePath"` // certificate storage path in blob storage
	Thumbprint        kmsdoc.HexStringStroable  `json:"thumbprint"`

	Template ResourceLocator `json:"template"` // locator for certificate template doc
	Issuer   ResourceLocator `json:"issuer"`   // locator for certificate doc for the actual issuer certificate
}

func (d *CertDoc) populateRef(dst *models.CertificateRefComposed) bool {
	if ok := d.BaseDoc.PopulateResourceRef(&dst.ResourceRef); !ok {
		return ok
	}
	dst.SubjectCommonName = d.SubjectCommonName
	dst.Thumbprint = d.Thumbprint.String()
	dst.NotAfter = d.NotAfter.Time()
	dst.Template = d.Template
	dst.Thumbprint = d.Thumbprint.String()
	return true
}

func (d *CertDoc) toModelRef() (r *models.CertificateRefComposed) {
	r = new(models.CertificateRefComposed)
	d.populateRef(r)
	return
}

func (d *CertDoc) toModel() *models.CertificateInfoComposed {
	r := new(models.CertificateInfoComposed)
	d.populateRef(&r.CertificateRefComposed)
	r.Issuer = d.Issuer
	d.KeySpec.PopulateKeyProperties(&r.Jwk)
	r.NotBefore = d.NotBefore.Time()
	r.Usages = d.Usages
	return r
}
