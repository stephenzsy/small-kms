package admin

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type adminServer struct {
	azBlobClient          *azblob.Client
	azBlobContainerClient *azblobcontainer.Client
}

func (s *adminServer) msGraphAppClient() (*msgraph.GraphServiceClient, error) {
	return nil, nil
}
