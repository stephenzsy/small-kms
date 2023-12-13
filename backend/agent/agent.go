package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	agentconfigmanager "github.com/stephenzsy/small-kms/backend/agent/configmanager/v2"
	agentendpoint "github.com/stephenzsy/small-kms/backend/agent/endpoint"

	agentpush "github.com/stephenzsy/small-kms/backend/agent/push"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

const DefaultEnvVarTenantID = "AZURE_TENANT_ID"
const DefaultEnvVarClientID = "AZURE_CLIENT_ID"
const DefaultEnvVarCertBundlePath = "AZURE_CERT_BUNDLE"
const DefaultEnvVarApiBaseUrl = "SMALLKMS_API_BASE_URL"
const DefaultEnvVarApiScope = "SMALLKMS_API_SCOPE"

var BuildID = "dev"

func main() {
	// Find .env file
	envFilePathPtr := flag.String("env", "", "path to .env file")
	envPrettyLog := flag.Bool("pretty-log", false, "pretty log")
	// slotPtr := flag.Uint("slot", 0, "slot")
	// //	skipTlsPtr := flag.Bool("skip-tls", false, "skip tls")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
	c := logger.WithContext(context.Background())

	if *envFilePathPtr != "" {
		logger.Printf("Loading environment file: %s\n", *envFilePathPtr)
		err := godotenv.Load(*envFilePathPtr)
		if err != nil {
			log.Printf("Error loading environment file: %s: %s\n", *envFilePathPtr, err.Error())
		}
	}

	envSvc := common.NewEnvService()

	agentPushEndpoint := envSvc.Default(agentcommon.EnvKeyAgentPushEndpoint, "https://localhost:8443", common.IdentityEnvVarPrefixAgent)

	args := flag.Args()

	slot := agentcommon.AgentSlot(args[0])
	if len(args) >= 2 {
		switch slot {
		case agentcommon.AgentSlotPrimary,
			agentcommon.AgentSlotSecondary,
			agentcommon.LegacyAgentSlotPrimary,
			agentcommon.LegacyAgentSlotSecondary:
			cm2, err := agentconfigmanager.NewConfigManager(envSvc, slot)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to create config manager")
			}
			cm2Poller := agentconfigmanager.NewConfigManagerPollTaskExecutor(cm2)

			// configManager, err := keeper.NewConfigManager(envSvc, slot)
			// if err != nil {
			// 	logger.Fatal().Err(err).Msg("Failed to create config manager")
			// }
			dockerClient, err := agentcommon.GetDockerClient()
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to create docker client")
			}

			// radiusConfigProcessor := radius.NewRadiusConfigProcessHandler(configManager.EnvConfig, configmanager.NewConfigDir(string(managedapp.AgentConfigNameRadius), configManager.ConfigDir))
			// configHandlerChain := &configmanager.ChainedContextConfigHandler{
			// 	ContextConfigHandler: radiusConfigProcessor,
			// }
			// radiusLauncher, err := radius.NewRadiusContainerLauncher(dockerClient, envSvc)
			// if err != nil {
			// 	logger.Fatal().Err(err).Msg("failed to create radius container launcher")
			// }
			// configHandlerChain.SetNext(&configmanager.ChainedContextConfigHandler{
			// 	ContextConfigHandler: radiusLauncher,
			// })
			// radiusConfigManager := radius.NewRadiusConfigManager(configHandlerChain, configManager.EnvConfig, configManager.ConfigDir)

			agentPushServer, err := agentpush.NewServer(BuildID, slot, envSvc, dockerClient)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to create agent push server")
			}
			newEcho := func(config *agentconfigmanager.AgentEndpointConfiguration) (*echo.Echo, error) {
				var err error
				e := echo.New()
				e.Use(middleware.Logger())
				e.Use(middleware.Recover())
				e.Use(ctx.InjectServiceContextMiddleware(context.Background()))
				e.Use(auth.PreconfiguredKeysJWTAuthorization(config.VerifyJWKs, agentPushEndpoint))
				e.Use(base.HandleResponseError)
				agentpush.RegisterHandlers(e, agentPushServer)
				agentendpoint.RegisterHandlers(e, agentPushServer)
				agentPushServer.SetEndpointConfig(config)

				e.TLSServer.Addr = args[1]
				e.TLSServer.TLSConfig, err = agentconfigmanager.GetTLSDefaultConfig(cm2)
				if err != nil {
					return nil, err
				}
				return e, nil
			}
			echoTask := agentconfigmanager.NewEchoTask(BuildID, newEcho, cm2, agentPushEndpoint, slot)

			tm := taskmanager.NewChainedTaskManager().
				WithTask(taskmanager.IntervalExecutorTask(cm2Poller, 0)).
				WithTask(echoTask)
			logger.Fatal().Err(taskmanager.StartWithGracefulShutdown(c, tm)).Msg("task manager exited")
			return
		}
	}

	fmt.Printf("Version: %s\n", BuildID)
	fmt.Printf("Usage: %s server 8443\n", os.Args[0])
}
