package main

import (
	"errors"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stephenzsy/small-kms/backend/agent/agentserver"
)

const EnvVarNameTLSCertPath = "AGENT_TLS_CERT_PATH"
const EnvVarNameTLSKeyPath = "AGENT_TLS_KEY_PATH"

func getCerts() (string, string, error) {
	certPath := os.Getenv(EnvVarNameTLSCertPath)
	keyPath := os.Getenv(EnvVarNameTLSKeyPath)
	if certPath == "" || keyPath == "" {
		return "", "", errors.New("missing TLS cert or key path")
	}
	return certPath, keyPath, nil
}

func bootstrapServer(addr string, skipTLS bool) {

	certPath, keyPath, err := getCerts()
	if err != nil && !skipTLS {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	server := agentserver.NewServer()

	agentserver.RegisterHandlers(e, server)
	if skipTLS {
		e.Logger.Fatal(e.Start(addr))
	} else {
		e.Logger.Fatal(e.StartTLS(addr, certPath, keyPath))
	}

}
