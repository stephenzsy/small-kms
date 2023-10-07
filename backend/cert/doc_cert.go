package cert

import (
	"crypto/x509"
	"fmt"

	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertificateStatus string

const (
	CertStatusInitialized CertificateStatus = "initialized"
	CertStatusPending     CertificateStatus = "pending"
	CertStatusIssued      CertificateStatus = "issued"
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
	PendingExpires    *kmsdoc.TimeStorable      `json:"pendingExpires"` // pending status expires time

	Template ResourceLocator `json:"template"` // locator for certificate template doc
	Issuer   ResourceLocator `json:"issuer"`   // locator for certificate doc for the actual issuer certificate
}

// PopulateX509 implements CertificateFieldsProvider.
func (doc *CertDoc) PopulateX509(cert *x509.Certificate) error {
	if doc.Status != CertStatusInitialized && doc.Status != CertStatusPending {
		return fmt.Errorf("certficiate doc status error: %s", doc.Status)
	}
	cert.SerialNumber = doc.SerialNumber.BigInt()
	cert.Subject.CommonName = doc.SubjectCommonName
	cert.NotBefore = doc.NotBefore.Time()
	cert.NotAfter = doc.NotAfter.Time()
	usageSet := utils.NewSet(doc.Usages...)
	if usageSet.Contains(models.CertUsageCA) {
		cert.IsCA = true
		if !usageSet.Contains(models.CertUsageCARoot) {
			cert.MaxPathLen = 1
		} else {
			cert.MaxPathLenZero = true
		}
		cert.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
	} else {
		if usageSet.Contains(models.CertUsageClientAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
		if usageSet.Contains(models.CertUsageServerAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
	}
	return nil
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

var _ CertificateFieldsProvider = (*CertDoc)(nil)
