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
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	requestcontext "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/managedapp"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/profile"
)

var BuildID = "dev"

func main() {
	if len(os.Args) < 3 {
		log.Info().Msg("Usage: smallkms <role> <listenerAddress>")
		os.Exit(1)
	}

	// Find .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}

	log.Printf("Server started, version: %s\n", BuildID)
	role := os.Args[1]
	listenerAddress := os.Args[2]
	if len(listenerAddress) == 0 {
		log.Error().Msg("listernerAddress is required")
		os.Exit(1)
	}

	switch role {
	case "admin":
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		if os.Getenv("ENABLE_CORS") == "true" {
			e.Use(middleware.CORS())
		}
		ctx := context.Background()
		server := api.NewServer()
		apiServer, err := api.NewApiServer(ctx, BuildID, server)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize api server")
		}
		e.Use(base.HandleResponseError)
		e.Use(requestcontext.InjectServiceContextMiddleware(apiServer))
		if os.Getenv("ENABLE_DEV_AUTH") == "true" {
			e.Use(auth.UnverifiedAADJwtAuth)
		} else {
			e.Use(auth.ProxiedAADAuth)
		}
		e.Use(server.GetAfterAuthMiddleware())
		models.RegisterHandlers(e, server)
		base.RegisterHandlers(e, base.NewBaseServer(BuildID))
		profile.RegisterHandlers(e, profile.NewServer(apiServer))
		managedapp.RegisterHandlers(e, managedapp.NewServer(apiServer))
		cert.RegisterHandlers(e, cert.NewServer(apiServer))
		//key.RegisterHandlers(e, key.NewServer(apiServer))
		common.StartEchoWithGracefulShutdown(ctx, e, func(ee *echo.Echo, shutdownNotifier common.LeafShutdownNotifier) {
			defer func() {
				shutdownNotifier.MarkShutdownComplete()
			}()
			log.Info().Err(ee.Start(listenerAddress)).Msg("echo server stopped")
		}, time.Minute)
	}

}
