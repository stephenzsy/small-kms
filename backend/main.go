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
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/admin"
	adminserver "github.com/stephenzsy/small-kms/backend/admin/server"
	agentpush "github.com/stephenzsy/small-kms/backend/agent/push"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	requestcontext "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/key"
	"github.com/stephenzsy/small-kms/backend/managedapp"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/secret"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
	"golang.org/x/net/http2"
)

var BuildID = "dev"

func main() {

	// Find .env file
	envFilePathPtr := flag.String("env-file", "", "path to .env file")
	envPrettyLog := flag.Bool("pretty-log", false, "pretty log")
	envDebug := flag.Bool("debug", false, "log")
	flag.Parse()

	if *envDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logger := log.Logger
	if *envPrettyLog {
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		logger = zerolog.New(output).With().Timestamp().Logger()

	}

	args := flag.Args()

	if len(args) < 2 {
		log.Info().Msg("Usage: smallkms <role> <listenerAddress>")
		os.Exit(1)
	}

	// Find .env file
	if *envFilePathPtr != "" {
		err := godotenv.Load(*envFilePathPtr)
		if err != nil {
			log.Printf("Error loading .env file: %s\n", err)
		}
	}

	logger.Info().Msgf("Server starting, version: %s", BuildID)
	role := args[0]
	listenerAddress := args[1]
	if len(listenerAddress) == 0 {
		log.Error().Msg("listernerAddress is required")
		os.Exit(1)
	}

	switch role {
	case "admin":
		e := echo.New()
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogMethod: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info().
					Str("URI", v.URI).
					Str("method", v.Method).
					Int("status", v.Status).
					Msg("request")
				return nil
			},
		}))
		e.Use(middleware.Recover())
		if os.Getenv("ENABLE_CORS") == "true" {
			e.Use(middleware.CORS())
		}
		ctx := logger.WithContext(context.Background())
		apiServer, err := api.NewApiServer(ctx, BuildID)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initialize api server")
		}
		e.Use(base.HandleResponseError)
		e.Use(requestcontext.InjectServiceContextMiddleware(apiServer))
		if os.Getenv("ENABLE_DEV_AUTH") == "true" {
			e.Use(auth.UnverifiedAADJwtAuth)
		} else {
			e.Use(auth.ProxiedAADAuth)
		}
		profile.RegisterHandlers(e, profile.NewServer(apiServer))
		managedapp.RegisterHandlers(e, managedapp.NewServer(apiServer))
		cert.RegisterHandlers(e, cert.NewServer(apiServer))
		agentpush.RegisterHandlers(e, agentpush.NewProxiedServer(apiServer))
		secret.RegisterHandlers(e, secret.NewServer(apiServer))
		key.RegisterHandlers(e, key.NewServer(apiServer))
		adminServer, err := adminserver.NewServer(apiServer)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initialize admin server")
		}
		admin.RegisterHandlers(e, adminServer)

		tm := taskmanager.NewChainedTaskManager().WithTask(
			taskmanager.NewTask("echo", func(c context.Context, sigCh <-chan os.Signal) error {
				if os.Getenv("ENABLE_DEV_AUTH") == "true" {
					go e.StartTLS(listenerAddress, "cert.pem", "key.pem")
				} else {
					s := &http2.Server{
						MaxConcurrentStreams: 250,
						MaxReadFrameSize:     1048576,
						IdleTimeout:          10 * time.Second,
					}
					go e.StartH2CServer(listenerAddress, s)
				}
				<-sigCh
				return e.Shutdown(c)
			}))
		logger.Fatal().Err(taskmanager.StartWithGracefulShutdown(ctx, tm)).Msg("task manager exited")
	}

}
