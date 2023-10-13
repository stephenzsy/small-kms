package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	server, err := agentserver.NewServer()
	if err != nil {
		panic(err)
	}

	agentserver.RegisterHandlers(e, server)

	mu := sync.RWMutex{}

	c := context.Background()
	loadConfigCh := make(chan string)
	httpReadyCh := make(chan bool)
	canShutdown := false
	go func() {
		for {
			switch {
			case <-httpReadyCh:
				mu.RLock()
			}
			canShutdown = true
			err := e.Start(addr)
			canShutdown = false
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info().Msg("server closed, reloading")
				} else {
					log.Error().Err(err).Msg("server error unexpected")
				}
			}

			mu.RUnlock()
		}
	}()
	go server.ConfigLoader.Start(c, loadConfigCh)

	for {
		select {
		case configVersion := <-loadConfigCh:
			log.Info().Msgf("config version changed: %s", configVersion)
			if canShutdown {
				err := e.Shutdown(c)
				if err != nil {
					log.Error().Err(err).Msg("failed to shutdown server")
				}
			}
			mu.Lock()
			httpReadyCh <- true
			mu.Unlock()

		}
	}
}
