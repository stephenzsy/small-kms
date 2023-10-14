package cm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var (
	ErrGracefullyShutdown   = errors.New("gracefully shutdown")
	ErrForcedShutdown       = errors.New("forced shutdown")
	ErrConfigProcessorError = errors.New("config processor error")
)

type configMsgType int

const (
	configMsgTypePoll configMsgType = iota + 1
	configMsgTypePush
)

type ConfigProcessor interface {
	StartPoll(ctx context.Context, scheduleToUpdate chan<- configMsg) error
	Process(ctx context.Context, msg configMsg) (bool, error)
}

type configMsg struct {
	msgType configMsgType
	name    shared.AgentConfigName
	data    shared.AgentConfiguration
}

type ConfigManager struct {
	isActive                 bool
	processors               map[shared.AgentConfigName]ConfigProcessor
	updateCh                 chan configMsg
	activeProcessorCtxCancel *context.CancelFunc
}

func (m *ConfigManager) cancelCurrentActiveProcessor(onDone chan<- struct{}) {
	if m.activeProcessorCtxCancel != nil {
		(*m.activeProcessorCtxCancel)()
	}
	onDone <- struct{}{}
}

// PushConfig implements ConfigManager.
func (m *ConfigManager) PushConfig(c context.Context, config shared.AgentConfiguration) {
	m.updateCh <- configMsg{
		msgType: configMsgTypePush,
		data:    config,
	}
}

// Shutdown implements ConfigManager.
func (m *ConfigManager) Shutdown(c context.Context) error {
	m.isActive = false
	onCanceledGracefully := make(chan struct{}, 1)
	go m.cancelCurrentActiveProcessor(onCanceledGracefully)
	for {
		select {
		case <-onCanceledGracefully:
			return nil
		case <-c.Done():
			return fmt.Errorf("%w:%w", ErrForcedShutdown, c.Err())
		}
	}
}

func shutdownServer(ctx context.Context, serverLaunched *echo.Echo) error {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return serverLaunched.Shutdown(c)
}

// Start implements ConfigManager.
func (m *ConfigManager) Start(c context.Context, onLaunchServer func(context.Context) (*echo.Echo, error)) error {
	m.updateCh = make(chan configMsg, len(m.processors)*3)
	for _, v := range m.processors {
		err := v.StartPoll(c, m.updateCh)
		if err != nil {
			return err
		}
	}
	m.isActive = true
	var serverStarted *echo.Echo
	for m.isActive {
		select {
		case <-c.Done():
			m.isActive = false
			return fmt.Errorf("%w:%w", ErrForcedShutdown, c.Err())
		case msg := <-m.updateCh:
			processor, ok := m.processors[msg.name]
			if !ok {
				log.Error().Err(ErrConfigProcessorError).Msgf("config processor not found: %s", msg.name)
				continue
			}
			hasChange, err := processor.Process(c, msg)
			if err != nil {
				log.Error().Err(err).Msg("config processor error")
				continue
			}
			if msg.name == shared.AgentConfigNameActiveServer {
				if serverStarted != nil && hasChange {
					err = shutdownServer(c, serverStarted)
					serverStarted = nil
					if err != nil {
						log.Error().Err(err).Msg("shutdown server error")
						continue
					}
				}
				if serverStarted == nil {
					serverStarted, err = onLaunchServer(c)
					if err != nil {
						log.Error().Err(err).Msg("launch server error")
						continue
					}
				}
			}
		}
	}
	return ErrGracefullyShutdown
}

func NewConfigManager() (*ConfigManager, error) {
	return &ConfigManager{}, nil
}
