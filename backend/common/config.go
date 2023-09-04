package common

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

type ServerRole string

const (
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
)

type commonConfig struct {
	defaultAzCerdential  *azidentity.DefaultAzureCredential
	keyvaultEndpoint     string
	azKeysClient         *azkeys.Client
	azCertificatesClient *azcertificates.Client
}

func NewCommonConfig() (c commonConfig, err error) {
	c.defaultAzCerdential, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return
	}
	c.keyvaultEndpoint = MustGetenv(DefualtEnvVarAzKeyvaultResourceEndpoint)
	c.azKeysClient, err = azkeys.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	if err != nil {
		return
	}
	c.azCertificatesClient, err = azcertificates.NewClient(c.keyvaultEndpoint, c.defaultAzCerdential, nil)
	return
}

type CommonConfig interface {
	DefaultAzCredential() azcore.TokenCredential
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
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
