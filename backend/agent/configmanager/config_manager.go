package cm

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
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
	Start(ctx context.Context, pollUpdate chan<- pollConfigMsg, exitCh chan<- error)
	Shutdown()
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
	isActive     bool
	processors   map[shared.AgentConfigName]ConfigProcessor
	sharedConfig sharedConfig
}

func (m *ConfigManager) Shutdown() {
	m.isActive = false
	for _, v := range m.processors {
		v.Shutdown()
	}
}

func (m *ConfigManager) Start(c context.Context, cmExitCh chan<- error) {
	pollUpdate := make(chan pollConfigMsg, 1)
	exitCh := make(chan error, len(m.processors))
	activeCount := 0
	for _, v := range m.processors {
		go v.Start(c, pollUpdate, exitCh)
		activeCount++
	}
	m.isActive = true
	for m.isActive {
		select {
		case <-c.Done():
			cmExitCh <- ErrForcedShutdown
			return
		case msg := <-pollUpdate:
			p := msg.processor
			err := p.Process(c, msg.task)
			if err != nil {
				log.Error().Err(err).Msgf("config processor error: %s", p.ConfigName())
			}
			p.MarkProcessDone(msg.task, err)
		case err := <-exitCh:
			activeCount--
			if err != nil {
				log.Error().Err(err).Msg("config processor exited")
			} else {
				log.Info().Err(err).Msg("config processor exited")
			}
		}
	}
	for activeCount > 0 {
		select {
		case <-c.Done():
			cmExitCh <- ErrForcedShutdown
			return
		case err := <-exitCh:
			activeCount--
			if err != nil {
				log.Error().Err(err).Msg("config processor exited")
			} else {
				log.Info().Err(err).Msg("config processor exited")
			}
		}
	}
	cmExitCh <- ErrGracefullyShutdown
}

func NewConfigManager(buildID string, configDir string) (*ConfigManager, error) {
	m := ConfigManager{
		processors: make(map[shared.AgentConfigName]ConfigProcessor, 1),
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
	m.processors[shared.AgentConfigNameActiveServer] = newActiveServerProcessor(&m.sharedConfig)
	return &m, nil
}

func StartConfigManagerWithGracefulShutdown(c context.Context, m *ConfigManager) {
	c = log.Logger.WithContext(c)
	c, cancel := context.WithCancel(c)
	defer cancel()
	exitErrCh := make(chan error, 1)

	go m.Start(c, exitErrCh)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	toCtx, toCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer toCancel()
	m.Shutdown()
	select {
	case <-toCtx.Done():
		cancel()
		log.Ctx(c).Fatal().Err(toCtx.Err()).Msg("failed to gracefully shutdown")
	case <-quit:
		cancel()
		log.Ctx(c).Fatal().Msg("forced shutdown")
	case err := <-exitErrCh:
		log.Ctx(c).Info().Err(err).Msg("gracefully shutdown")
	}
}
