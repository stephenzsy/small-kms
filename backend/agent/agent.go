package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/agent/keeper"
	"github.com/stephenzsy/small-kms/backend/agent/taskmanager"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
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

	configDir := common.MustGetenv("AGENT_CONFIG_DIR")
	configManager, err := keeper.NewConfigManager(configDir)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create config manager")
	}

	agentPort := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "GATEWAY_PORT", "8443")
	agentPortNum, err := strconv.Atoi(agentPort)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to parse AGENT_GATEWAY_PORT")
	}

	args := flag.Args()
	if len(args) >= 2 {
		switch args[0] {
		case "server":
			newEcho := func(config keeper.AgentServerConfiguration) (*echo.Echo, error) {
				var err error
				e := echo.New()
				e.Use(middleware.Logger())
				e.Use(middleware.Recover())
				e.Use(base.HandleResponseError)
				base.RegisterHandlers(e, base.NewBaseServer(BuildID))

				e.TLSServer.Addr = args[1]
				e.TLSServer.TLSConfig, err = keeper.GetTLSDefaultConfig(config)
				if err != nil {
					return nil, err
				}
				return e, nil
			}
			keeperTask := keeper.NewKeeper(configManager)
			echoTask := keeper.NewEchoTask(BuildID, newEcho, keeperTask, agentPortNum)

			tm := taskmanager.NewChainedTaskManager().
				WithTask(taskmanager.IntervalExecutorTask(keeperTask, 0)).
				WithTask(echoTask)
			logger.Fatal().Err(taskmanager.StartWithGracefulShutdown(c, tm)).Msg("task manager exited")
			// e := echo.New()
			// e.Use(middleware.Logger())
			// e.Use(middleware.Recover())
			// e.TLSServer.Addr = args[1]
			// configManager, err := cm.NewConfigManager(BuildID, configDir, args[1], uint32(*slotPtr))
			// if err != nil {
			// 	log.Panicf("Failed to create config manager: %v\n", err)
			// }
			// s, err := agentserver.NewServer()
			// if err != nil {
			// 	log.Panicf("Failed to create server: %v\n", err)
			// }
			// agentserver.RegisterHandlers(e, s)

			// configManager.Manage(e)

			// cm.StartConfigManagerWithGracefulShutdown(context.Background(), configManager)
			// //bootstrapServer(args[1], *skipTlsPtr)
			return
		}
	}

	fmt.Printf("Version: %s\n", BuildID)
	fmt.Printf("Usage: %s server 8443\n", os.Args[0])
}
