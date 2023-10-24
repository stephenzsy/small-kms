package cert

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertificateStatus string

const (
	CertStatusInitialized CertificateStatus = "initialized"
	CertStatusPending     CertificateStatus = "pending"
	CertStatusIssued      CertificateStatus = "issued"
)

type CertJwkSpec struct {
	ct.CertKeySpec
	KID     string                            `json:"kid"`
	X5u     *string                           `json:"x5u,omitempty"`
	X5t     shared.Base64RawURLEncodableBytes `json:"x5t,omitempty"`
	X5tS256 shared.Base64RawURLEncodableBytes `json:"x5t#S256,omitempty"`

	keyExportable bool
}

const queryColumnTemplate = "c.template"
const queryColumnStatus = "c.status"
const queryColumnNotAfter = "c.notAfter"
const queryColumnThumbprint = "c.thumbprint"
const (
	CertDocQueryColumnCertStorePath = "c.certStorePath"
	CertDocQueryColumnNotBefore     = "c.notBefore"
)

type CertDoc struct {
	kmsdoc.BaseDoc

	Status CertificateStatus `json:"status"` // certificate status

	// X509 certificate info
	SerialNumber      SerialNumberStorable            `json:"serialNumber"`
	SubjectCommonName string                          `json:"subjectCommonName"`
	NotBefore         kmsdoc.TimeStorable             `json:"notBefore"`
	NotAfter          kmsdoc.TimeStorable             `json:"notAfter"`
	Usages            []shared.CertificateUsage       `json:"usages"`
	CertSpec          CertJwkSpec                     `json:"certSpec"`
	KeyStorePath      *string                         `json:"keyStorePath,omitempty"`
	CertStorePath     string                          `json:"certStorePath"` // certificate storage path in blob storage
	Thumbprint        shared.CertificateFingerprint   `json:"thumbprint"`
	PendingExpires    *kmsdoc.TimeStorable            `json:"pendingExpires"` // pending status expires time
	TemplateDigest    kmsdoc.HexStringStroable        `json:"templateDigest"` // copied from template doc
	Template          shared.ResourceLocator          `json:"template"`       // locator for certificate template doc
	Issuer            shared.ResourceLocator          `json:"issuer"`         // locator for certificate doc for the actual issuer certificate
	SANs              *shared.SubjectAlternativeNames `json:"sans,omitempty"` // subject alternative names
}

// SnapshotWithNewLocator implements kmsdoc.KmsDocumentSnapshotable.
func (doc *CertDoc) SnapshotWithNewLocator(locator shared.ResourceLocator) *CertDoc {
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
	Thumbprint    shared.CertificateFingerprint
	Issuer        shared.ResourceLocator
}

func (doc *CertDoc) FetchCertificatePEMBlob(c context.Context) ([]byte, error) {
	blobClient := common.GetAdminServerClientProvider(c).CertsAzBlobContainerClient()
	get, err := blobClient.NewBlobClient(doc.CertStorePath).DownloadStream(c, nil)
	if err != nil {
		return nil, err
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(c, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		return nil, err
	}

	err = retryReader.Close()
	if err != nil {
		return nil, err

	}
	return downloadedData.Bytes(), nil
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
	if usageSet.Contains(shared.CertUsageCA) {
		cert.BasicConstraintsValid = true
		cert.IsCA = true
		if usageSet.Contains(shared.CertUsageCARoot) {
			cert.MaxPathLen = 1
		} else {
			cert.MaxPathLenZero = true
		}
		cert.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
	} else {
		cert.KeyUsage |= x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		if usageSet.Contains(shared.CertUsageClientAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
		if usageSet.Contains(shared.CertUsageServerAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
	}
	if doc.SANs != nil {
		cert.DNSNames = doc.SANs.DNSNames
		cert.IPAddresses = doc.SANs.IPAddresses
		cert.EmailAddresses = doc.SANs.Emails
	}
	return &cert, nil
}

func (d *CertDoc) populateRef(r *shared.CertificateRef) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.SubjectCommonName = d.SubjectCommonName
	r.Thumbprint = d.Thumbprint
	r.NotAfter = d.NotAfter.Time()
	r.Template = d.Template
	r.IsIssued = d.Status == CertStatusIssued
}

func (d *CertDoc) toModelRef() (r *shared.CertificateRef) {
	if d == nil {
		return nil
	}
	r = new(shared.CertificateRef)
	d.populateRef(r)
	return
}

func (d *CertDoc) toModel() *shared.CertificateInfo {
	if d == nil {
		return nil
	}
	r := new(shared.CertificateInfo)
	d.populateRef(&r.CertificateRef)
	r.Issuer = d.Issuer
	d.CertSpec.PopulateKeyProperties(&r.Jwk)
	r.NotBefore = d.NotBefore.Time()
	r.Usages = d.Usages
	r.SubjectAlternativeNames = d.SANs
	return r
}

func (k *CertJwkSpec) PopulateKeyProperties(r *shared.JwkProperties) {
	if k == nil || r == nil {
		return
	}
	k.CertKeySpec.PopulateKeyProperties(r)
	r.KeyID = &k.KID
	r.CertificateURL = k.X5u
	r.CertificateThumbprint = k.X5t
	r.CertificateThumbprintSHA256 = k.X5tS256
}

var _ kmsdoc.KmsDocumentSnapshotable[*CertDoc] = (*CertDoc)(nil)
