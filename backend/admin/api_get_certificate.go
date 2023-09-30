package admin

import (
	"bytes"
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
)

// deprecated
func (s *adminServer) FetchCertificatePEMBlob(ctx context.Context, blobName string) ([]byte, error) {
	blobClient := s.azBlobContainerClient
	get, err := blobClient.NewBlobClient(blobName).DownloadStream(ctx, nil)
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

func (c *CertificateRef) populateFromDoc(certDoc *CertDoc) {
	c.Deleted = certDoc.Deleted
	c.ID = certDoc.ID.GetUUID()
	c.Issuer = certDoc.IssuerCertificateID.GetUUID()
	c.IssuerNamespace = certDoc.IssuerNamespaceID
	c.Name = certDoc.CommonName
	c.NamespaceID = certDoc.NamespaceID
	c.NotAfter = certDoc.NotAfter
	c.Updated = certDoc.Updated
	c.UpdatedBy = certDoc.UpdatedBy
	c.Usage = certDoc.Usage
}

func (s *adminServer) GetCertificateV1(c *gin.Context, namespaceID uuid.UUID, id uuid.UUID, params GetCertificateV1Params) {
	if _, ok := authNamespaceRead(c, namespaceID); !ok {
		return
	}

	certIdentifier := CertificateIdentifier{ID: id, Type: params.ByType}
	certDoc, err := s.getCertDoc(c, namespaceID, certIdentifier.docID())
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		log.Error().Err(err).Msg("Internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	cert := new(CertificateRef)
	cert.populateFromDoc(certDoc)
	if params.Format == nil || (*params.Format != FormatPEM && *params.Format != FormatJWK) {
		c.JSON(http.StatusOK, cert)
		return
	}
	// fetchPemBlob
	if *params.Format == FormatPEM {
		c.JSON(http.StatusOK, cert)
		return
	}
	c.JSON(http.StatusOK, cert)
}
