/*
 * Small KMS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
)

func main() {
	log.Println("Server started")

	log.Printf("Server role: %s", os.Getenv("APP_ROLE"))
	if os.Getenv("APP_ROLE") == "admin" || false {
		// Perform DB migration
		// check if file exists
		if _, err := os.Stat("/app/data/smallkms.db"); err != nil {
			// file does not ecist
			f, err := os.Open("/app/data/smallkms.db")
			if err != nil {
				log.Panicf("Failed to create DB file: %s", err.Error())
			}
			err = f.Close()
			if err != nil {
				log.Panicf("Failed to create DB file: %s", err.Error())
			}
		}
		m, err := migrate.New(
			"file:///app/migrations",
			"sqlite3:///app/data/smallkms.db")
		m.Steps(2)
		if err != nil {
			log.Panicf("Failed to perform DB migration: %s", err.Error())
		}
		m.Close()
	}
	serverConfig := common.NewServerConfig()
	router := gin.Default()

	handleAadAuth := func(c *gin.Context) {
		// Intercept the headers here
		authedPrincipalId := c.Request.Header.Get("X-Ms-Client-Principal-Id")
		if !serverConfig.IsPrincipalIdTrusted(authedPrincipalId) {
			c.JSON(403, gin.H{"error": fmt.Sprintf("Principal ID is not trusted: %s", authedPrincipalId)})
			c.Abort()
			return
		}

		c.Next()
	}

	router.Use(handleAadAuth)

	switch serverConfig.GetServerRole() {
	case common.ServerRoleAdmin:
		admin.RegisterHandlers(router, admin.NewAdminServer(&serverConfig))
	}

	log.Fatal(router.Run(":9000"))
}
