package admin

import (
	"bytes"
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *adminServer) GetCertificateV1(c *gin.Context, namespaceID uuid.UUID, id uuid.UUID, params GetCertificateV1Params) {

	c.JSON(http.StatusNotFound, nil)
}
