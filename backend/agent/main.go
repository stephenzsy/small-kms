package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stephenzsy/small-kms/backend/agent/agentserver"
	cm "github.com/stephenzsy/small-kms/backend/agent/configmanager"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

const DefaultEnvVarTenantID = "AZURE_TENANT_ID"
const DefaultEnvVarClientID = "AZURE_CLIENT_ID"
const DefaultEnvVarCertBundlePath = "AZURE_CERT_BUNDLE"
const DefaultEnvVarApiBaseUrl = "SMALLKMS_API_BASE_URL"
const DefaultEnvVarApiScope = "SMALLKMS_API_SCOPE"

var myNamespacIdentifier = shared.StringIdentifier("me")

func _oldmain() {
	// Find .env file
	envFilePathPtr := flag.String("env", "", "path to .env file")
	slotPtr := flag.Uint("slot", 0, "slot")
	//	skipTlsPtr := flag.Bool("skip-tls", false, "skip tls")
	flag.Parse()

	if *envFilePathPtr != "" {
		err := godotenv.Load(*envFilePathPtr)
		if err != nil {
			log.Printf("Error loading environment file: %s: %s\n", *envFilePathPtr, err.Error())
		}
	}

	configDir := common.MustGetenv("AGENT_CONFIG_DIR")

	args := flag.Args()
	if len(args) >= 2 {
		switch args[0] {
		case "bootstrap":
			switch args[1] {
			case "active-host":
				return
			}
		case "server":
			e := echo.New()
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			e.TLSServer.Addr = args[1]
			configManager, err := cm.NewConfigManager(BuildID, configDir, args[1], uint32(*slotPtr))
			if err != nil {
				log.Panicf("Failed to create config manager: %v\n", err)
			}
			s, err := agentserver.NewServer()
			if err != nil {
				log.Panicf("Failed to create server: %v\n", err)
			}
			agentserver.RegisterHandlers(e, s)

			configManager.Manage(e)

			cm.StartConfigManagerWithGracefulShutdown(context.Background(), configManager)
			//bootstrapServer(args[1], *skipTlsPtr)
			return
		}
	}

	fmt.Printf("Version: %s\n", BuildID)
	fmt.Printf("Usage: %s bootstrap active-host\n", os.Args[0])
	fmt.Printf("Usage: %s server :8443\n", os.Args[0])

}
