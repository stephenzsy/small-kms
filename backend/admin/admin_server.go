package admin

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/common"
)

type adminServer struct {
	common.CommonConfig
	azBlobClient                 *azblob.Client
	azBlobContainerClient        *azblobcontainer.Client
	azCosmosClient               *azcosmos.Client
	azCosmosDatabaseClient       *azcosmos.DatabaseClient
	azCosmosContainerClientCerts *azcosmos.ContainerClient
	msGraphClient                *msgraph.GraphServiceClient
	skipDeviceCheck              bool
}

func NewAdminServer() *adminServer {
	cosmosEndpoint := common.MustGetenv(common.DefualtEnvVarAzCosmosResourceEndpoint)
	storageBlobEndpoint := common.MustGetenv(common.DefualtEnvVarAzStroageBlobResourceEndpoint)

	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic(err)
	}
	s := adminServer{
		CommonConfig:    &commonConfig,
		skipDeviceCheck: os.Getenv("SKIP_DEVICE_CHECK") == "true",
	}
	s.azBlobClient, err = azblob.NewClient(storageBlobEndpoint, s.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get az blob client: %s", err.Error())
	}
	s.azBlobContainerClient = s.azBlobClient.ServiceClient().NewContainerClient(common.GetEnvWithDefault("AZURE_STORAGEBLOB_CONTAINERNAME_CERTS", "certs"))

	s.azCosmosClient, err = azcosmos.NewClient(cosmosEndpoint, s.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get az cosmos client: %s", err.Error())
	}
	s.azCosmosDatabaseClient, err = s.azCosmosClient.NewDatabase(common.GetEnvWithDefault("AZURE_COSMOS_DATABASE_ID", "kms"))
	if err != nil {
		log.Panicf("Failed to get az cosmos database client: %s", err.Error())
	}
	s.azCosmosContainerClientCerts, err = s.azCosmosDatabaseClient.NewContainer(common.GetEnvWithDefault("AZURE_COSMOS_CONTAINERNAME_CERTS", "Certs"))
	if err != nil {
		log.Panicf("Failed to get az cosmos container client for Certs: %s", err.Error())
	}

	s.msGraphClient, err = msgraph.NewGraphServiceClientWithCredentials(s.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get graph clients: %s", err.Error())
	}
	return &s
}
