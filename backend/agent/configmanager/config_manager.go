package cm

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var (
	ErrGracefullyShutdown   = errors.New("gracefully shutdown")
	ErrForcedShutdown       = errors.New("forced shutdown")
	ErrConfigProcessorError = errors.New("config processor error")
)

type ConfigProcessor interface {
	Name() shared.AgentConfigName
	Start(ctx context.Context, pollUpdate chan<- pollConfigMsg, exitCh chan<- error)
	Shutdown()
	Process(ctx context.Context) error
	MarkProcessDone()
}

type pollConfigMsg struct {
	name      shared.AgentConfigName
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
			err := p.Process(c)
			if err != nil {
				log.Error().Err(err).Msgf("config processor error: %s", p.Name())
			}
			p.MarkProcessDone()
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

func NewConfigManager(buildID string) (*ConfigManager, error) {
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
	} else {
		m.sharedConfig.init(buildID, agentClient)
	}

	m.processors[shared.AgentConfigNameHeartbeat] = newHeartbeatConfigProcessor(&m.sharedConfig)
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
