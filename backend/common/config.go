package common

import (
	"database/sql"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
)

type ServerRole string

const (
	ServerRoleAdmin ServerRole = "admin"
)

type ServerConfig interface {
	GetServerRole() ServerRole
	GetDB() *sql.DB
	GetAzKeysClient() *azkeys.Client
	GetAzBlobClient() *azblob.Client
	GetAzBlobContainerName() string
}

type serverConfig struct {
	role                  string
	azKeyVaultEndpoint    string
	azBlobServiceEndpoint string
	azBlobContainerName   string

	db           *sql.DB
	azCredential *azidentity.DefaultAzureCredential
	azKeysClient *azkeys.Client
	azBlobClient *azblob.Client
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

	if err := config.initDB(); err != nil {
		log.Panicf("Failed to initialize DB: %s", err.Error())
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

	return config
}

func (config *serverConfig) initDB() (err error) {
	log.Println("Initialize DB")
	config.db, err = sql.Open("sqlite3", "data/smallkms.db")
	return
}
