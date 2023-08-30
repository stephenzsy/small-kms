package common

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type ServerRole string

const (
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzCosmosClientID              = "AZURE_COSMOS_CLIENTID"
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzKeyvaultClientID            = "AZURE_KEYVAULT_CLIENTID"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobClientID         = "AZURE_STORAGEBLOB_CLIENTID"
)

func GetAzCredential(clientId string) (azcore.TokenCredential, error) {
	if len(clientId) > 0 {
		return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(clientId),
		})
	}
	return azidentity.NewDefaultAzureCredential(nil)
}

func GetAzCosmosClient() (*azcosmos.Client, error) {
	credential, err := GetAzCredential(os.Getenv(DefualtEnvVarAzCosmosClientID))
	if err != nil {
		return nil, err
	}
	return azcosmos.NewClient(os.Getenv(DefualtEnvVarAzCosmosResourceEndpoint), credential, nil)
}

func GetAzKeysClient() (*azkeys.Client, error) {
	credential, err := GetAzCredential(os.Getenv(DefualtEnvVarAzKeyvaultClientID))
	if err != nil {
		return nil, err
	}
	return azkeys.NewClient(os.Getenv(DefualtEnvVarAzKeyvaultResourceEndpoint), credential, nil)
}

func GetAzStorageBlobClient() (*azblob.Client, error) {
	credential, err := GetAzCredential(os.Getenv(DefualtEnvVarAzStroageBlobClientID))
	if err != nil {
		return nil, err
	}
	return azblob.NewClient(os.Getenv(DefualtEnvVarAzStroageBlobResourceEndpoint), credential, nil)
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
