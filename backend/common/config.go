package common

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type ServerRole string

const (
	DefaultEnvVarAzureTenantId                 = "AZURE_TENANT_ID"
	DefaultEnvVarAzureClientId                 = "AZURE_CLIENT_ID"
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
	DefaultEnvVarAzureManagedIdentityClientId  = "AZURE_MANAGED_IDENTITY_CLIENT_ID"
)

type commonConfig struct {
	defaultAzCerdential          azcore.TokenCredential
	keyvaultEndpoint             string
	azKeysClient                 *azkeys.Client
	azCertificatesClient         *azcertificates.Client
	azCosmosClient               *azcosmos.Client
	azCosmosDatabaseClient       *azcosmos.DatabaseClient
	azCosmosContainerClientCerts *azcosmos.ContainerClient
	msGraphClient                *msgraph.GraphServiceClient
	tenantID                     uuid.UUID
}

func NewCommonConfig() (c commonConfig, err error) {
	if managedIdentityClientId, ok := os.LookupEnv(DefaultEnvVarAzureManagedIdentityClientId); ok {
		c.defaultAzCerdential, err = azidentity.NewManagedIdentityCredential(
			&azidentity.ManagedIdentityCredentialOptions{
				ID: azidentity.ClientID(managedIdentityClientId),
			})
	} else {
		c.defaultAzCerdential, err = azidentity.NewDefaultAzureCredential(nil)
	}
	if err != nil {
		return
	}
	c.keyvaultEndpoint = MustGetenv(DefualtEnvVarAzKeyvaultResourceEndpoint)
	c.azKeysClient, err = azkeys.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	if err != nil {
		return
	}
	c.azCertificatesClient, err = azcertificates.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	c.tenantID = uuid.MustParse(MustGetenv(DefaultEnvVarAzureTenantId))

	cosmosEndpoint := MustGetenv(DefualtEnvVarAzCosmosResourceEndpoint)
	c.azCosmosClient, err = azcosmos.NewClient(cosmosEndpoint, c.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get az cosmos client: %s", err.Error())
	}
	c.azCosmosDatabaseClient, err = c.azCosmosClient.NewDatabase(GetEnvWithDefault("AZURE_COSMOS_DATABASE_ID", "kms"))
	if err != nil {
		log.Panicf("Failed to get az cosmos database client: %s", err.Error())
	}
	c.azCosmosContainerClientCerts, err = c.azCosmosDatabaseClient.NewContainer(GetEnvWithDefault("AZURE_COSMOS_CONTAINERNAME_CERTS", "Certs"))
	if err != nil {
		log.Panicf("Failed to get az cosmos container client for Certs: %s", err.Error())
	}

	c.msGraphClient, err = msgraph.NewGraphServiceClientWithCredentials(c.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to get graph clients: %s", err.Error())
	}

	return
}

type CommonConfig interface {
	DefaultAzCredential() azcore.TokenCredential
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
	AzCosmosContainerClient() *azcosmos.ContainerClient
	MsGraphClient() *msgraph.GraphServiceClient
	TenantID() uuid.UUID
}

func (c *commonConfig) DefaultAzCredential() azcore.TokenCredential {
	return c.defaultAzCerdential
}

func (c *commonConfig) AzKeysClient() *azkeys.Client {
	return c.azKeysClient
}

func (c *commonConfig) AzCertificatesClient() *azcertificates.Client {
	return c.azCertificatesClient
}

func (c *commonConfig) AzCosmosContainerClient() *azcosmos.ContainerClient {
	return c.azCosmosContainerClientCerts
}

func (c *commonConfig) MsGraphClient() *msgraph.GraphServiceClient {
	return c.msGraphClient
}

func (c *commonConfig) TenantID() uuid.UUID {
	return c.tenantID
}

func MustGetenv(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = %s", name, value)
	return
}

func GetEnvWithDefault(name string, defaultValue string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		value = defaultValue
	}
	log.Printf("Config %s = %s", name, value)
	return
}
