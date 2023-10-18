package cm

import (
	"context"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
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
	Start(ctx context.Context, pollUpdate chan<- pollConfigMsg, shutdownNotifier common.LeafShutdownNotifier)
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
}

func (m *ConfigManager) Start(c context.Context, shutdownNotifier common.LeafShutdownNotifier) {
	isActive := true
	pollUpdate := make(chan pollConfigMsg, 1)
	shutdownNotifiers := make([]common.LeafShutdownNotifier, 0, len(m.processors))
	for _, v := range m.processors {
		notifier := common.NewLeafShutdownNotifier(string(v.ConfigName()))
		shutdownNotifiers = append(shutdownNotifiers, notifier)
		go v.Start(c, pollUpdate, notifier)
	}
	mergedNotifier := common.MergeShutdownNotifier("config manager processors", shutdownNotifiers...)
	for isActive {
		select {
		case sig := <-shutdownNotifier.Quit():
			isActive = false
			mergedNotifier.RelaySingal(sig)
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
	shutdownNotifier.RelayComplete(mergedNotifier)
}

func NewConfigManager(buildID string, configDir string) (*ConfigManager, error) {
	m := ConfigManager{
		processors: make(map[shared.AgentConfigName]ConfigProcessor, 2),
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

	m.processors[shared.AgentConfigNameHeartbeat] = newHeartbeatConfigProcessor(&m.sharedConfig)
	m.serverReadyCh = make(chan ActiveServerReadyConfig, 1)
	m.processors[shared.AgentConfigNameActiveServer] = newActiveServerProcessor(&m.sharedConfig, m.serverReadyCh)
	return &m, nil
}

func (m *ConfigManager) Manage(e *echo.Echo) {
	m.managedServer = e
}

func (m *ConfigManager) startManage(ctx context.Context, shutdownNotifier common.LeafShutdownNotifier) {
	defer shutdownNotifier.MarkShutdownComplete()
	isActive := true
	logger := log.Ctx(ctx)
	hasStarted := false
	e := m.managedServer
	for isActive {
		select {
		case <-shutdownNotifier.Quit():
			isActive = false
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
			tlsConfig.MinVersion = tls.VersionTLS12
			allowedCertFingerprints := make(map[[sha512.Size384]byte]bool, len(msg.AuthorizedClientCertificateFingerprintsSHA384))
			for _, v := range msg.AuthorizedClientCertificateFingerprintsSHA384 {
				allowedCertFingerprints[[sha512.Size384]byte(v)] = true
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
	shutdownNotifier := common.NewLeafShutdownNotifier("main")
	echoShutdownNotifier := common.NewLeafShutdownNotifier("echo server")
	go m.startManage(c, echoShutdownNotifier)
	cmShutdownNotifier := common.NewLeafShutdownNotifier("config manager")
	go m.Start(c, cmShutdownNotifier)

	mainShutdownNotifier := common.MergeShutdownNotifier("master", shutdownNotifier, echoShutdownNotifier, cmShutdownNotifier)

	mainShutdownNotifier.ListenOSIntercept()
	<-shutdownNotifier.Quit()
	shutdownNotifier.MarkShutdownComplete()
	toCtx, toCancel := context.WithTimeout(c, 10*time.Second)
	defer toCancel()
	m.managedServer.Shutdown(toCtx)
	select {
	case <-toCtx.Done():
		log.Ctx(c).Fatal().Err(toCtx.Err()).Msg("failed to gracefully shutdown")
	case <-shutdownNotifier.Quit():
		log.Ctx(c).Fatal().Msg("forced shutdown after second interupt")
	case <-mainShutdownNotifier.Complete():
		log.Ctx(c).Info().Msg("gracefully shutdown")
	}
}
