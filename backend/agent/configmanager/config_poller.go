package configmanager

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type ConfigPoller[CK comparable, T any] struct {
	ChainedContextConfigHandler
	name         string
	contextKey   CK
	onPollConfig func(context.Context) (T, time.Duration, error)
}

func jitter(d time.Duration) time.Duration {
	jitterInterval := time.Minute * 5
	if d >= time.Hour {
		jitterInterval = time.Minute
	} else {
		jitterInterval = time.Second * 5
	}
	return d + time.Duration(rand.Int63n(int64(jitterInterval)))
}

func (p *ConfigPoller[CK, T]) Start(c context.Context, exit <-chan os.Signal) error {
	logger := log.Ctx(c).With().Str("configPoller", p.name).Logger()
	logger.Debug().Msg("starting config poller")
	defer logger.Debug().Msg("config poller stopped")

	timer := time.NewTimer(0)
	active := true
	for active {
		select {
		case <-c.Done():
			return c.Err()
		case <-exit:
			active = false
			timer.Stop()
		case <-timer.C:
			timer.Stop()
			polledConfig, nextDuration, err := p.onPollConfig(c)
			if err != nil {
				logger.Error().Err(err).Msg("error while polling config")
			}
			if !active {
				return nil
			}
			if _, err := p.Handle(context.WithValue(c, p.contextKey, polledConfig)); err != nil {
				logger.Error().Err(err).Msg("error while handling config")
			}
			if nextDuration < 5*time.Minute {
				nextDuration = 5 * time.Minute
			}
			nextDuration = jitter(nextDuration)
			logger.Debug().Dur("nextPoolInMillis", nextDuration).Msg("next config poll")
			timer.Reset(nextDuration)
		}
	}
	return nil
}

func (p *ConfigPoller[CK, T]) ReceivePushedConfig(c context.Context, config T) (context.Context, error) {
	return p.Handle(context.WithValue(c, p.contextKey, config))
}

func (p *ConfigPoller[CK, T]) Name() string {
	return p.name
}

func NewConfigPoller[CK comparable, T any](name string, contextKey CK, onPollConfig func(context.Context) (T, time.Duration, error)) *ConfigPoller[CK, T] {
	return &ConfigPoller[CK, T]{
		name:         name,
		contextKey:   contextKey,
		onPollConfig: onPollConfig,
	}
}

var _ taskmanager.Task = (*ConfigPoller[string, any])(nil)
