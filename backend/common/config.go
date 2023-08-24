package common

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerRole string

const (
	ServerRoleAdmin ServerRole = "admin"
)

type ServerConfig interface {
	GetServerRole() ServerRole
	GetDB() *sql.DB
}

type serverConfig struct {
	role string
	db   *sql.DB
}

func (c *serverConfig) GetDB() *sql.DB {
	return c.db
}

func (c *serverConfig) GetServerRole() ServerRole {
	return ServerRole(c.role)
}

func NewServerConfig() serverConfig {
	config := serverConfig{}
	godotenv.Load(".env")

	config.role = os.Getenv("APP_ROLE")
	if config.role == "" {
		log.Panicln("No APP_ROLE configured")
	}
	switch config.role {
	case string(ServerRoleAdmin):
		break
	default:
		log.Panicf("Unknown APP_ROLE: %s", config.role)
	}

	if err := config.initDB(); err != nil {
		log.Panicf("Failed to initialize DB: %s", err.Error())
	}

	return config
}

func (config *serverConfig) initDB() (err error) {
	log.Println("Initialize DB")
	config.db, err = sql.Open("sqlite3", "data/smallkms.db")
	return
}
