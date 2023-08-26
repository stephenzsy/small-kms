package common

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type ServerRole string

const (
	ServerRoleAdmin ServerRole = "admin"
)

type ServerConfig interface {
	GetServerRole() ServerRole
	// Deprecated: migrating
	GetDB() *sql.DB
	GetAzKeysClient() *azkeys.Client
	GetAzBlobClient() *azblob.Client
	GetAzBlobContainerName() string
	AzCosmosContainerClient() *azcosmos.ContainerClient
	IsPrincipalIdTrusted(principalId string) bool
}

type serverConfig struct {
	role                     string
	azKeyVaultEndpoint       string
	azBlobServiceEndpoint    string
	azBlobContainerName      string
	azCosmosEndpoint         string
	azCosmosDatabaseId       string
	azCosmosCertsContainerId string

	db                      *sql.DB
	azCredential            *azidentity.DefaultAzureCredential
	azKeysClient            *azkeys.Client
	azBlobClient            *azblob.Client
	azCosmosClient          *azcosmos.Client
	azCosmosDbClient        *azcosmos.DatabaseClient
	azCosmosContainerClient *azcosmos.ContainerClient
	trustedPrincipalIds     map[string]bool
	skipTrustFilter         bool
}

func (c *serverConfig) GetDB() *sql.DB {
	return c.db
}

func (c *serverConfig) GetServerRole() ServerRole {
	return ServerRole(c.role)
}

func (c *serverConfig) GetAzKeysClient() *azkeys.Client {
	return c.azKeysClient
}

func (c *serverConfig) GetAzBlobClient() *azblob.Client {
	return c.azBlobClient
}

func (c *serverConfig) GetAzBlobContainerName() string {
	return c.azBlobContainerName
}

func (c *serverConfig) AzCosmosContainerClient() *azcosmos.ContainerClient {
	return c.azCosmosContainerClient
}

func (c *serverConfig) IsPrincipalIdTrusted(principalId string) bool {
	if c.skipTrustFilter {
		return true
	}
	if len(principalId) == 0 {
		return false
	}
	return c.trustedPrincipalIds[principalId]
}

func mustGetenv(name string) (value string) {
	value = os.Getenv(name)
	if len(value) == 0 {
		log.Panicf("No variable %s configured", name)
	}
	log.Printf("Config %s = %s", name, value)
	return
}

func NewServerConfig() serverConfig {
	config := serverConfig{}
	godotenv.Load(".env")

	config.role = mustGetenv("APP_ROLE")
	switch config.role {
	case string(ServerRoleAdmin):
		break
	default:
		log.Panicf("Unknown APP_ROLE: %s", config.role)
	}

	var err error = nil
	if config.azCredential, err = azidentity.NewDefaultAzureCredential(nil); err != nil {
		log.Panicf("Failed to initialize azure credential: %s", err.Error())
	}
	config.azKeyVaultEndpoint = mustGetenv("AZURE_KEY_VAULT_ENDPOINT")
	if config.azKeysClient, err = azkeys.NewClient(config.azKeyVaultEndpoint, config.azCredential, nil); err != nil {
		log.Panicf("Failed to initialize key vault client: %s", err.Error())
	}
	config.azBlobServiceEndpoint = mustGetenv("AZURE_BLOB_SERVICE_ENDPOINT")
	if config.azBlobClient, err = azblob.NewClient(config.azBlobServiceEndpoint, config.azCredential, nil); err != nil {
		log.Panicf("Failed to initialize blob client: %s", err.Error())
	}
	config.azBlobContainerName = mustGetenv("AZURE_BLOB_CONTAINER_NAME")
	config.azCosmosEndpoint = mustGetenv("AZURE_COSMOS_ENDPOINT")
	if config.azCosmosClient, err = azcosmos.NewClient(config.azCosmosEndpoint, config.azCredential, nil); err != nil {
		log.Panicf("Failed to initialize cosmos client: %s", err.Error())
	}
	config.azCosmosDatabaseId = mustGetenv("AZURE_COSMOS_DATABASE_ID")
	if config.azCosmosDbClient, err = config.azCosmosClient.NewDatabase(config.azCosmosDatabaseId); err != nil {
		log.Panicf("Failed to initialize cosmos database client: %s", err.Error())
	}
	config.azCosmosCertsContainerId = os.Getenv("AZURE_COSMOS_CERTS_CONTAINER_ID")
	if len(config.azCosmosCertsContainerId) == 0 {
		config.azCosmosCertsContainerId = "Certs"
	}
	if config.azCosmosContainerClient, err = config.azCosmosDbClient.NewContainer(config.azCosmosCertsContainerId); err != nil {
		log.Panicf("Failed to initialize cosmos container client: %s", err.Error())
	}

	config.trustedPrincipalIds = make(map[string]bool)
	trustedIdsEnv := os.Getenv("TRUSTED_SERVICE_PRINCIPAL_IDS")
	if trustedIdsEnv == "disabled" {
		config.skipTrustFilter = true
	} else {
		for _, s := range strings.Split(trustedIdsEnv, ",") {
			if len(s) > 0 {
				parsed, err := uuid.Parse(s)
				if err != nil {
					log.Panicf("Failed to parse trusted principal ID: %s", s)

				}
				config.trustedPrincipalIds[parsed.String()] = true
			}
		}
	}

	return config
}
