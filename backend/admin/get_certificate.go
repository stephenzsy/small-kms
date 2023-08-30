package admin

import (
	"bytes"
	"context"
	"encoding/pem"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *adminServer) GetCertificateV1(c *gin.Context, namespaceID NamespaceID, id uuid.UUID, params GetCertificateV1Params) {
	result, err := s.ReadCertDBItem(c, namespaceID, id)
	if err != nil {
		log.Printf("Faild to get certificate metadata: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if result.ID == uuid.Nil {
		log.Printf("Faild to get certificate metadata: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	accept := AcceptJson
	switch *params.Accept {
	case AcceptX509CaCert:
		accept = AcceptX509CaCert
	case AcceptPem:
		accept = AcceptPem
	}

	if accept == AcceptJson {
		c.JSON(200, result.CertificateRef)
		return
	}

	if len(result.CertStore) == 0 {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	pemBlob, err := s.FetchCertificatePEMBlob(c, result.CertStore)
	if err != nil {
		log.Printf("Faild to fetch certificate blob: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	switch *params.Accept {
	case AcceptX509CaCert:
		block, _ := pem.Decode(pemBlob)
		if block == nil {
			log.Printf("Faild to decode certificate blob stored")
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		c.Data(200, "application/x-x509-ca-cert", block.Bytes)
	default:
		c.Data(200, "application/x-pem-file", pemBlob)
	}
}
