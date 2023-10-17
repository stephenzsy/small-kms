package cm

import (
	"context"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var (
	ErrGracefullyShutdown = errors.New("gracefully shutdown")
	ErrForcedShutdown     = errors.New("forced shutdown")
	ErrSaveConfigFile     = errors.New("config file save error")
)

type ConfigProcessor interface {
	ConfigName() shared.AgentConfigName
	Version() string
	Start(ctx context.Context, pollUpdate chan<- pollConfigMsg)
	Process(ctx context.Context, task string) error
	MarkProcessDone(taskName string, err error)
}

type pollConfigMsg struct {
	name      shared.AgentConfigName
	task      string
	processor ConfigProcessor
}

type ConfigManager struct {
	common.CommonServer
	processors    map[shared.AgentConfigName]ConfigProcessor
	sharedConfig  sharedConfig
	serverReadyCh chan ActiveServerReadyConfig
	managedServer *echo.Echo
	shutdownCtrl  *ShutdownController
}

func (m *ConfigManager) Start(c context.Context) {
	shutdownCh, shutdownDefer := m.shutdownCtrl.Subsribe()
	defer shutdownDefer()
	pollUpdate := make(chan pollConfigMsg, 1)
	for _, v := range m.processors {
		go v.Start(c, pollUpdate)
	}
	for m.shutdownCtrl.IsActive() {
		select {
		case <-shutdownCh:
		case <-c.Done():
			log.Error().Err(c.Err()).Msg("config manager forced shutdown")
			return
		case msg := <-pollUpdate:
			p := msg.processor
			err := p.Process(c, msg.task)
			if err != nil {
				log.Error().Err(err).Msgf("config processor error: %s", p.ConfigName())
			}
			p.MarkProcessDone(msg.task, err)
		}
	}
}

func NewConfigManager(buildID string, configDir string) (*ConfigManager, error) {
	m := ConfigManager{
		processors:   make(map[shared.AgentConfigName]ConfigProcessor, 2),
		shutdownCtrl: NewShutdownController(),
	}
	var err error
	if m.CommonServer, err = common.NewCommonConfig(); err != nil {
		return nil, err
	}
	serviceIdentity := m.CommonServer.ServiceIdentity()
	if apiBasePath, err := common.GetNonEmptyEnv("API_BASE_PATH"); err != nil {
		return nil, err
	} else if apiScope, err := common.GetNonEmptyEnv("API_SCOPE"); err != nil {
		return nil, err
	} else if agentClient, err := agentclient.NewClientWithCreds(apiBasePath, serviceIdentity.TokenCredential(), []string{apiScope}, serviceIdentity.TenantID()); err != nil {
		return nil, err
	} else if azKeyVaultUrl, err := common.GetNonEmptyEnv("AZURE_KEYVAULT_RESOURCEENDPOINT"); err != nil {
		return nil, err
	} else if azSecretsClient, err := azsecrets.NewClient(azKeyVaultUrl, serviceIdentity.TokenCredential(), nil); err != nil {
		return nil, err
	} else {
		m.sharedConfig.init(buildID, agentClient, azSecretsClient, configDir)
	}

	m.processors[shared.AgentConfigNameHeartbeat] = newHeartbeatConfigProcessor(&m.sharedConfig, m.shutdownCtrl)
	m.serverReadyCh = make(chan ActiveServerReadyConfig, 1)
	m.processors[shared.AgentConfigNameActiveServer] = newActiveServerProcessor(&m.sharedConfig, m.serverReadyCh, m.shutdownCtrl)
	return &m, nil
}

func (m *ConfigManager) Manage(e *echo.Echo) {
	m.managedServer = e
}

func (m *ConfigManager) startManage(ctx context.Context) {
	shutDownCh, shutdownDefer := m.shutdownCtrl.Subsribe()
	defer shutdownDefer()
	logger := log.Ctx(ctx)
	hasStarted := false
	e := m.managedServer
	for m.shutdownCtrl.IsActive() {
		select {
		case <-shutDownCh:
		case msg := <-m.serverReadyCh:
			if hasStarted {
				err := e.Shutdown(ctx)
				if err != nil {
					logger.Error().Err(err).Msg("failed to shutdown server")
				}
				hasStarted = false
			}
			tlsConfig := new(tls.Config)
			certContent, err := os.ReadFile(msg.ServerCertificateFile)
			if err != nil {
				logger.Error().Err(err).Msg("failed to read server certificate file")
				continue
			}
			tlsCert, err := tls.X509KeyPair(certContent, certContent)
			if err != nil {
				logger.Error().Err(err).Msg("failed to parse server certificate file")
				continue
			}
			tlsConfig.Certificates = []tls.Certificate{tlsCert}
			tlsConfig.ClientCAs = x509.NewCertPool()
			_, rest := pem.Decode(certContent)
			_, rest = pem.Decode(rest)
			tlsConfig.ClientCAs.AppendCertsFromPEM(rest)
			tlsConfig.ClientAuth = tls.RequireAnyClientCert
			allowedCertFingerprints := make(map[[sha512.Size384]byte]bool, len(msg.AuthorizedClientCertificateFingerprintsSHA384))
			for _, v := range msg.AuthorizedClientCertificateFingerprintsSHA384 {
				allowedCertFingerprints[v] = true
			}
			tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
				for _, rawCert := range rawCerts {
					if allowedCertFingerprints[sha512.Sum384(rawCert)] {
						return nil
					}
				}
				return &tls.CertificateVerificationError{}
			}
			e.TLSServer.TLSConfig = tlsConfig
			hasStarted = true
			log.Info().Err(e.StartServer(e.TLSServer)).Msg("server exited")
		}
	}
}

func StartConfigManagerWithGracefulShutdown(c context.Context, m *ConfigManager) {
	c = log.Logger.WithContext(c)
	c, cancel := context.WithCancel(c)
	defer cancel()
	go m.startManage(c)
	go m.Start(c)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	toCtx, toCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer toCancel()
	m.shutdownCtrl.Shutdown()
	m.managedServer.Shutdown(toCtx)
	select {
	case <-toCtx.Done():
		cancel()
		log.Ctx(c).Fatal().Err(toCtx.Err()).Msg("failed to gracefully shutdown")
	case <-quit:
		cancel()
		log.Ctx(c).Fatal().Msg("forced shutdown")
	case <-m.shutdownCtrl.C:
		log.Ctx(c).Info().Msg("gracefully shutdown")
	}
}
