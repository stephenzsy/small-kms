package common

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
)

type ServerRole string

const (
	DefaultEnvVarAzureTenantId                 = "AZURE_TENANT_ID"
	DefaultEnvVarAzureClientId                 = "AZURE_CLIENT_ID"
	DefaultEnvVarAppAzureClientId              = "APP_AZURE_CLIENT_ID"
	DefaultEnvVarAppAzureClientSecret          = "APP_AZURE_CLIENT_SECRET"
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
	DefaultEnvVarAzureManagedIdentityClientId  = "AZURE_MANAGED_IDENTITY_CLIENT_ID"
	DefualtEnvVarAzSubscriptionID              = "AZURE_SUBSCRIPTION_ID"
	DefaultEnvVarAzResourceGroupName           = "AZURE_RESOURCE_GROUP_NAME"
)

type commonConfig struct {
	defaultAzCerdential          azcore.TokenCredential
	azKeyvaultName               string
	keyvaultEndpoint             string
	azKeysClient                 *azkeys.Client
	azCertificatesClient         *azcertificates.Client
	azCosmosClient               *azcosmos.Client
	azCosmosDatabaseClient       *azcosmos.DatabaseClient
	azCosmosContainerClientCerts *azcosmos.ContainerClient
	tenantIDStr                  string
	tenantID                     uuid.UUID
	aadAppClientId               string
	aadAppClientSecret           string
	confidentialAppCredential    azcore.TokenCredential
	subscriptionId               string
	azResourceGroupName          string
}

// AzKeyVaultName implements CommonConfig.
func (s *commonConfig) AzKeyvaultName() string {
	return s.azKeyvaultName
}

// AzResourceGroupName implements CommonConfig.
func (c *commonConfig) AzResourceGroupName() string {
	return c.azResourceGroupName
}

// SubscriptionID implements CommonConfig.
func (s *commonConfig) AzSubscriptionID() string {
	return s.subscriptionId
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
	{
		parsed, _ := url.Parse(c.keyvaultEndpoint)
		c.azKeyvaultName = strings.Split(parsed.Host, ".")[0]
		if len(c.azKeyvaultName) == 0 {
			log.Panicf("unable to parse keyvault name from key vault url")
		}
	}
	c.azKeysClient, err = azkeys.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	if err != nil {
		return
	}
	c.azCertificatesClient, err = azcertificates.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	if err != nil {
		return
	}
	c.tenantIDStr = MustGetenv(DefaultEnvVarAzureTenantId)
	c.tenantID = uuid.MustParse(MustGetenv(DefaultEnvVarAzureTenantId))
	c.tenantIDStr = c.tenantID.String()

	if cosmosConnStr, ok := os.LookupEnv("AZURE_COSMOS_CONNECTION_STRING"); ok {
		c.azCosmosClient, err = azcosmos.NewClientFromConnectionString(cosmosConnStr, nil)
	} else {
		cosmosEndpoint := MustGetenv(DefualtEnvVarAzCosmosResourceEndpoint)
		c.azCosmosClient, err = azcosmos.NewClient(cosmosEndpoint, c.DefaultAzCredential(), nil)
	}
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

	c.aadAppClientId = MustGetenv(DefaultEnvVarAppAzureClientId)
	c.aadAppClientSecret = MustGetenvSecret(DefaultEnvVarAppAzureClientSecret)
	c.subscriptionId = MustGetenv(DefualtEnvVarAzSubscriptionID)
	c.azResourceGroupName = MustGetenv(DefaultEnvVarAzResourceGroupName)

	c.confidentialAppCredential, err = azidentity.NewClientSecretCredential(c.tenantIDStr, c.aadAppClientId, c.aadAppClientSecret, nil)
	if err != nil {
		return
	}
	return
}

type CommonConfig interface {
	DefaultAzCredential() azcore.TokenCredential
	AzKeyvaultEndpoint() string
	AzKeyvaultName() string
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
	AzCosmosContainerClient() *azcosmos.ContainerClient
	TenantID() uuid.UUID
	AzSubscriptionID() string
	AzResourceGroupName() string
	ConfidentialAppCredential() azcore.TokenCredential
	NewOnBehalfOfCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (*azidentity.OnBehalfOfCredential, error)
}

func (c *commonConfig) DefaultAzCredential() azcore.TokenCredential {
	return c.defaultAzCerdential
}

func (c *commonConfig) AzKeyvaultEndpoint() string {
	return c.keyvaultEndpoint
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

func (c *commonConfig) TenantID() uuid.UUID {
	return c.tenantID
}

func (c *commonConfig) ConfidentialAppCredential() azcore.TokenCredential {
	return c.confidentialAppCredential
}

func (c *commonConfig) NewOnBehalfOfCredential(userAssertion string,
	opts *azidentity.OnBehalfOfCredentialOptions) (*azidentity.OnBehalfOfCredential, error) {
	return azidentity.NewOnBehalfOfCredentialWithSecret(c.tenantIDStr,
		c.aadAppClientId,
		userAssertion,
		c.aadAppClientSecret, opts)
}

var _ CommonConfig = (*commonConfig)(nil)

func MustGetenv(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = %s", name, value)
	return
}

func MustGetenvSecret(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = **********", name)
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
