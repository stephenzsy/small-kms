package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ready            bool
	cosmosEndpoint   string
	cosmosDatabaseID string
}

var config ServerConfig

func Initialize() {
	if config.ready {
		return
	}

	godotenv.Load(".env")

	config.cosmosEndpoint = os.Getenv("AZURE_COSMOS_ENDPOINT")
	if config.cosmosEndpoint == "" {
		log.Panicln("No cosmos DB endpoint configured")
	}
	log.Printf("Cosmos DB endpoint: %s", config.cosmosEndpoint)

	config.cosmosDatabaseID = os.Getenv("AZURE_COSMOS_DATABASE_ID")
	if config.cosmosDatabaseID == "" {
		log.Panicln("No cosmos DB database configured")
	}
	log.Printf("Cosmos DB database: %s", config.cosmosDatabaseID)

	if err := dbInit(); err != nil {
		log.Panicf("Failed to initialize DB: %s", err.Error())
	}

	config.ready = true
}

func GetCosmosDatabaseClient() (*azcosmos.DatabaseClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	client, err := azcosmos.NewClient(config.cosmosEndpoint, cred, nil)
	if err != nil {
		return nil, err
	}
	dbClient, err := client.NewDatabase(config.cosmosDatabaseID)
	if err != nil {
		return dbClient, err
	}

	// verify db exist
	if _, err := dbClient.Read(context.Background(), nil); err != nil {
		return dbClient, err
	}
	return dbClient, err
}

func dbInit() error {
	log.Println("Initialize DB")
	client, err := GetCosmosDatabaseClient()
	if err != nil {
		return err
	}
	// check and provision containers
	certsClient, err := client.NewContainer("Certificates")
	if err != nil {
		return err
	}
	resp, err := certsClient.Read(context.Background(), nil)
	if err != nil {
		return err
	}
	if len(resp.ContainerProperties.PartitionKeyDefinition.Paths) != 2 ||
		resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != "/Category" ||
		resp.ContainerProperties.PartitionKeyDefinition.Paths[1] != "/Name" {
		return fmt.Errorf("partition key not as spec: %s", resp.ContainerProperties.PartitionKeyDefinition.Paths)
	}
	return nil
}
