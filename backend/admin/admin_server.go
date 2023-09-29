package admin

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

	//"github.com/google/uuid"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	//msgraphapp "github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	//msgraphsp "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"

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
	//appClientId                  string
	//appServicePrincipalId        string
	// name to id mapping
	//appRoles map[string]uuid.UUID
}

func NewAdminServer() *adminServer {
	cosmosEndpoint := common.MustGetenv(common.DefualtEnvVarAzCosmosResourceEndpoint)
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
	/*
		s.appClientId = common.MustGetenv(common.DefaultEnvVarAzureClientId)
		spobj, err := s.msGraphClient.ServicePrincipalsWithAppId(&s.appClientId).Get(context.Background(),
			&msgraphsp.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
				QueryParameters: &msgraphsp.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
					Select: []string{"id"},
				},
			})
		if err != nil {
			log.Panicf("Failed to get current graph service principal: %s", err.Error())
		}
		s.appServicePrincipalId = *spobj.GetId()

		app, err := s.msGraphClient.ApplicationsWithAppId(&s.appClientId).Get(context.Background(),
			&msgraphapp.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
				QueryParameters: &msgraphapp.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
					Select: []string{"id", "appRoles"},
				},
			})
		if err != nil {
			log.Panicf("Failed to get current graph app: %s", err.Error())
		}
		appRoles := app.GetAppRoles()
		s.appRoles = make(map[string]uuid.UUID, len(appRoles))
		for _, appRole := range appRoles {
			s.appRoles[*appRole.GetValue()] = *appRole.GetId()
		}*/
	return &s
}
