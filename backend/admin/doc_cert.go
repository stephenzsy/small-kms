package admin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stephenzsy/small-kms/backend/common"
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
	doc := new(CertDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.AzCosmosContainerClient(), nsID, docID, doc)
	return doc, common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:cert:%s", nsID, docID))
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

func (doc *CertDoc) storeCertificatePEMBlob(ctx context.Context, blobClient *azblobcontainer.Client, b []byte) (*string, error) {
	blobName := doc.CertStorePath
	if blobName == "" {
		return nil, errors.New("empty blob name")
	}
	blockBlobClient := blobClient.NewBlockBlobClient(blobName)
	_, err := blockBlobClient.UploadBuffer(ctx, b, &blockblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr("application/x-pem-file"),
		},
		Metadata: map[string]*string{
			"issuer_id": to.Ptr(fmt.Sprintf("%s/%s", doc.IssuerNamespaceID, doc.IssuerCertificateID.GetUUID())),
			// "x5t":       base64UrlToHexStrPtr(doc.KeyInfo.CertificateThumbprint),
			// "x5t_S256":  base64UrlToHexStrPtr(doc.KeyInfo.CertificateThumbprintSHA256),
		},
	})
	return ToPtr(blockBlobClient.URL()), err
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
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	pager := s.AzCosmosContainerClient().NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.fingerprint FROM c
WHERE c.namespaceId = @namespaceId
  AND c.type = @type`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: nsID.String()},
				{Name: "@type", Value: kmsdoc.DocTypeNameCert},
			},
		})

	return PagerToList[CertDoc](ctx, pager)
}
