package admin

import (
	"bytes"
	"context"
	"encoding/pem"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
)

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
	c.Issuer = certDoc.IssuerID.GetUUID()
	c.IssuerNamespace = certDoc.IssuerNamespace
	c.Name = certDoc.Name
	c.NamespaceID = certDoc.NamespaceID
	c.NotAfter = certDoc.Expires
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
	pemBlob, err := s.FetchCertificatePEMBlob(c, certDoc.CertStorePath)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		log.Error().Err(err).Msg("Internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if *params.Format == FormatPEM {
		pemString := string(pemBlob)
		cert.Pem = &pemString
		c.JSON(http.StatusOK, cert)
		return
	}
	x5c := make([][]byte, 0)
	for block, rest := pem.Decode(pemBlob); block != nil; block, rest = pem.Decode(rest) {
		x5c = append(x5c, block.Bytes)
	}
	cert.X5c = &x5c
	c.JSON(http.StatusOK, cert)
}
