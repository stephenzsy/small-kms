package admin

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
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

func (s *adminServer) readCertDoc(ctx context.Context, nsID uuid.UUID, docID kmsdoc.KmsDocID) (*CertDoc, error) {
	return nil, nil
}

func (d *CertDoc) GetCUID() kmsdoc.KmsDocID {
	if d.AliasID != nil {
		return *d.AliasID
	}
	return d.ID
}

func (doc *CertDoc) fetchCertificatePEMBlob(ctx context.Context, blobClient *azblobcontainer.Client) ([]byte, error) {
	get, err := blobClient.NewBlobClient(doc.CertStorePath).DownloadStream(ctx, nil)
	if err != nil {
		return nil, err
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
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

func (s *adminServer) toCertificateInfo(ctx context.Context,
	doc *CertDoc,
	include *IncludeCertificateParameter,
	certPemBlob []byte) (*CertificateInfo, error) {
	if doc == nil {
		return nil, nil
	}
	certInfo := CertificateInfo{
		CommonName: doc.CommonName,
		Usage:      doc.Usage,
		NotBefore:  doc.NotBefore,
		NotAfter:   doc.NotAfter,
		Subject:    doc.Subject,
	}
	if doc.SubjectAlternativeNames != nil {
		certInfo.SubjectAlternativeNames = doc.SubjectAlternativeNames
	}

	baseDocPopulateRefWithMetadata(&doc.BaseDoc, &certInfo.Ref)
	docCuid := doc.GetCUID()
	certInfo.Ref.ID = docCuid.GetUUID()
	certInfo.Ref.Type = RefTypeCertificate
	certInfo.Ref.DisplayName = doc.FingerprintSHA1Hex

	if include != nil {
		if len(certPemBlob) == 0 {
			if fetchedCertPemBlob, err := doc.fetchCertificatePEMBlob(ctx, s.azBlobContainerClient); err != nil {
				return nil, err
			} else {
				certPemBlob = fetchedCertPemBlob
			}
		}
		switch *include {
		case IncludePEM:
			certInfo.Pem = ToPtr(string(certPemBlob))
		case IncludeJWK:
			// certInfo.Jwk.populateCertsFromPemBlob(certPemBlob)
		}
	}

	certInfo.Template.ID = doc.TemplateID.GetUUID()
	certInfo.Template.Type = RefTypeCertificateTemplate
	certInfo.Template.NamespaceID = doc.NamespaceID

	certInfo.IssuerCertificate.ID = doc.IssuerCertificateID.GetUUID()
	certInfo.IssuerCertificate.Type = RefTypeCertificate
	certInfo.IssuerCertificate.NamespaceID = doc.IssuerNamespaceID
	if doc.NamespaceID == doc.IssuerNamespaceID {
		// root CA
		certInfo.IssuerCertificate.ID = certInfo.Ref.ID
	}

	return &certInfo, nil
}

func (s *adminServer) listCertificateDocs(ctx context.Context, nsID uuid.UUID) ([]*CertDoc, error) {
	return nil, nil
}
