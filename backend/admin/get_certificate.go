package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *adminServer) GetCertificateV1(c *gin.Context, namespaceID NamespaceID, id uuid.UUID, params GetCertificateV1Params) {
	db := s.config.AzCosmosContainerClient()
	resp, err := db.ReadItem(c, azcosmos.NewPartitionKeyString(uuid.Nil.String()), id.String(), nil)
	if err != nil {
		log.Printf("Faild to get certificate metadata: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
	}
	result := CertDBItem{}
	err = json.Unmarshal(resp.Value, &result)
	if err != nil {
		log.Printf("Faild to unmarshall certificate metadata: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
	}
	if len(result.CertStore) == 0 {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	if params.Accept == nil || len(*params.Accept) == 0 || *params.Accept == AcceptJson {
		c.JSON(200, result.CertificateRef)
		return
	}

	filename := "cert.pem"
	contentType := AcceptPem
	switch *params.Accept {
	case AcceptX509CaCert:
		filename = "cert.der"
		contentType = AcceptX509CaCert
	}

	blobClient := s.config.GetAzBlobClient()

	get, err := blobClient.DownloadStream(c, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", result.CertStore, filename), nil)
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(c, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	err = retryReader.Close()
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.Data(200, string(contentType), downloadedData.Bytes())
}
