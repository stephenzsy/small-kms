package cert

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
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

type CertJwkSpec struct {
	ct.CertKeySpec
	KID     string                   `json:"kid"`
	X5t     kmsdoc.Base64UrlStorable `json:"x5t,omitempty"`
	X5tS256 kmsdoc.Base64UrlStorable `json:"x5t#S256,omitempty"`

	keyExportable bool
}

type CertDoc struct {
	kmsdoc.BaseDoc

	Status CertificateStatus `json:"status"` // certificate status

	// X509 certificate info
	SerialNumber      SerialNumberStorable      `json:"serialNumber"`
	SubjectCommonName string                    `json:"subjectCommonName"`
	NotBefore         kmsdoc.TimeStorable       `json:"notBefore"`
	NotAfter          kmsdoc.TimeStorable       `json:"notAfter"`
	Usages            []models.CertificateUsage `json:"usages"`
	CertSpec          CertJwkSpec               `json:"certSpec"`
	KeyStorePath      *string                   `json:"keyStorePath,omitempty"`
	CertStorePath     string                    `json:"certStorePath"` // certificate storage path in blob storage
	Thumbprint        kmsdoc.HexStringStroable  `json:"thumbprint"`
	PendingExpires    *kmsdoc.TimeStorable      `json:"pendingExpires"` // pending status expires time
	TemplateDigest    kmsdoc.HexStringStroable  `json:"templateDigest"` // copied from template doc
	Template          ResourceLocator           `json:"template"`       // locator for certificate template doc
	Issuer            ResourceLocator           `json:"issuer"`         // locator for certificate doc for the actual issuer certificate
}

// SnapshotWithNewLocator implements kmsdoc.KmsDocumentSnapshotable.
func (doc *CertDoc) SnapshotWithNewLocator(locator common.Locator[models.NamespaceKind, models.ResourceKind]) *CertDoc {
	if doc == nil {
		return nil
	}
	snapshotDoc := *doc

	snapshotDoc.BaseDoc.NamespaceID = locator.GetNamespaceID()
	snapshotDoc.BaseDoc.ID = locator.GetID()

	return &snapshotDoc
}

type CertDocSigningPatch struct {
	CertSpec      CertJwkSpec
	CertStorePath string
	Thumbprint    kmsdoc.HexStringStroable
	Issuer        ResourceLocator
}

func (d *CertDoc) patchSigned(c common.ServiceContext, patch *CertDocSigningPatch) error {
	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet("/thumbprint", patch.Thumbprint.HexString())
	patchOps.AppendSet("/certStorePath", patch.CertStorePath)
	patchOps.AppendSet("/issuer", patch.Issuer.String())
	patchOps.AppendSet("/status", CertStatusIssued)
	patchOps.AppendRemove("/pendingExpires")
	patchOps.AppendSet("/certSpec", patch.CertSpec)

	err := kmsdoc.Patch(c, d.GetLocator(), d, patchOps)
	if err != nil {
		return err
	}
	d.Thumbprint = patch.Thumbprint
	d.CertStorePath = patch.CertStorePath
	d.Issuer = patch.Issuer
	d.Status = CertStatusIssued
	d.PendingExpires = nil
	d.CertSpec = patch.CertSpec

	return nil
}

func createCertificateDoc(nsID models.NamespaceID,
	tmpl *ct.CertificateTemplateDoc,
	params models.IssueCertificateFromTemplateParams) (*CertDoc, error) {

	certID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	doc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID:   nsID,
			ID:            common.NewIdentifierWithKind(models.ResourceKindCert, common.UUIDIdentifier(certID)),
			SchemaVersion: 1,
		},
		Status:            CertStatusInitialized,
		SerialNumber:      certID[:],
		SubjectCommonName: tmpl.SubjectCommonName,
		Usages:            tmpl.Usages,
		CertSpec: CertJwkSpec{
			CertKeySpec: tmpl.KeySpec,
		},
		KeyStorePath:   tmpl.KeyStorePath,
		Template:       tmpl.GetLocator(),
		Issuer:         tmpl.IssuerTemplate,
		NotBefore:      kmsdoc.TimeStorable(now),
		NotAfter:       kmsdoc.TimeStorable(now.AddDate(0, int(tmpl.ValidityInMonths), 0)),
		TemplateDigest: tmpl.Digest,
	}

	return &doc, nil
}

func (doc *CertDoc) createX509Certificate() (*x509.Certificate, error) {
	if doc.Status != CertStatusInitialized && doc.Status != CertStatusPending {
		return nil, fmt.Errorf("certficiate doc status error: %s", doc.Status)
	}
	cert := x509.Certificate{}
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
	return &cert, nil
}

func (d *CertDoc) populateRef(r *models.CertificateRefComposed) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.SubjectCommonName = d.SubjectCommonName
	r.Thumbprint = d.Thumbprint.HexString()
	r.NotAfter = d.NotAfter.Time()
	r.Template = d.Template
}

func (d *CertDoc) toModelRef() (r *models.CertificateRefComposed) {
	if d == nil {
		return nil
	}
	r = new(models.CertificateRefComposed)
	d.populateRef(r)
	return
}

func (d *CertDoc) toModel() *models.CertificateInfoComposed {
	if d == nil {
		return nil
	}
	r := new(models.CertificateInfoComposed)
	d.populateRef(&r.CertificateRefComposed)
	r.Issuer = d.Issuer
	d.CertSpec.PopulateKeyProperties(&r.Jwk)
	r.NotBefore = d.NotBefore.Time()
	r.Usages = d.Usages
	return r
}

func (k *CertJwkSpec) PopulateKeyProperties(r *models.JwkProperties) {
	if k == nil || r == nil {
		return
	}
	k.CertKeySpec.PopulateKeyProperties(r)
	r.CertificateThumbprint = k.X5t.StringPtr()
	r.CertificateThumbprintSHA256 = k.X5tS256.StringPtr()
}

var _ kmsdoc.KmsDocumentSnapshotable[*CertDoc] = (*CertDoc)(nil)
