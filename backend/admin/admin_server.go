package admin

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
)

type adminServer struct {
	common.CommonConfig
	azBlobClient          *azblob.Client
	azBlobContainerClient *azblobcontainer.Client
	graphService          graph.GraphService
}

func NewAdminServer() *adminServer {
	storageBlobEndpoint := common.MustGetenv(common.DefualtEnvVarAzStroageBlobResourceEndpoint)

	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic(err)
	}
	s := adminServer{
		CommonConfig: &commonConfig,
	}
	s.azBlobClient, err = azblob.NewClient(storageBlobEndpoint, s.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get az blob client: %s", err.Error())
	}
	s.azBlobContainerClient = s.azBlobClient.ServiceClient().NewContainerClient(common.GetEnvWithDefault("AZURE_STORAGEBLOB_CONTAINERNAME_CERTS", "certs"))
	s.graphService = graph.NewGraphService(s.CommonConfig)

	return &s
}

func (s *adminServer) msGraphClient(c context.Context) (*msgraph.GraphServiceClient, error) {
	if authIdentity, ok := auth.GetAuthIdentity(c); ok {
		if creds, err := authIdentity.GetOnBehalfOfTokenCredential(s, nil); err != nil {
			return nil, err
		} else {
			return msgraph.NewGraphServiceClientWithCredentials(creds, nil)
		}
	}
	return nil, fmt.Errorf("%w: no auth header to authenticate to graph service", common.ErrStatusUnauthorized)
}

func (s *adminServer) msGraphAppClient() (*msgraph.GraphServiceClient, error) {
	return msgraph.NewGraphServiceClientWithCredentials(s.ConfidentialAppCredential(), nil)
}
