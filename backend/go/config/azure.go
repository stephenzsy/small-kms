package config

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/joho/godotenv"
)

var envReady bool = false

type ServiceConfig struct {
	cosmosEndpoint string
}

var config ServiceConfig

func ensureEnvReady() {
	if !envReady {
		godotenv.Load("../.env")
		envReady = true
	}

	config.cosmosEndpoint = os.Getenv("AZURE_COSMOS_ENDPOINT")
	if config.cosmosEndpoint == "" {
		panic("No cosmos db configured")
	}
}

func GetAzCosmosClient() (*azcosmos.Client, error) {
	ensureEnvReady()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return azcosmos.NewClient(config.cosmosEndpoint, cred, nil)
}
