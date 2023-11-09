package configmanager

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type VersionedConfig interface {
	GetVersion() string
	NextPullAfter() time.Time
}

type ConfigPoller[CK comparable, T VersionedConfig] struct {
	ChainedContextConfigHandler
	name             string
	contextKey       CK
	pullRemoteConfig func(context.Context) (T, error)
	cache            *ConfigCache[T]
}

func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
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

	// try load config from cache
	cachedConfig, cacheOK, err := p.cache.Load()
	if err != nil {
		logger.Error().Err(err).Msg("error while loading cached config")
	}
	firstDuration := time.Duration(0)
	if cacheOK {
		p.Handle(context.WithValue(c, p.contextKey, cachedConfig))
		firstDuration = time.Until(cachedConfig.NextPullAfter())
		if firstDuration < 5*time.Minute {
			firstDuration = 5 * time.Minute
		}
		firstDuration = jitter(firstDuration)
		logger.Debug().Dur("nextPoolInMillis", firstDuration).Msg("next config poll after cache load")
	}
	timer := time.NewTimer(firstDuration)
	active := true
	for active {
		select {
		case <-c.Done():
			return c.Err()
		case <-exit:
			active = false
		case <-timer.C:
			timer.Stop()
			pulled, err := p.pullRemoteConfig(c)
			if err != nil {
				logger.Error().Err(err).Msg("error while polling config")
			}
			if err := p.cache.SetPulledConfig(pulled); err != nil {
				logger.Error().Err(err).Msg("error while setting pulled config")
			}
			if active {
				if _, err := p.Handle(context.WithValue(c, p.contextKey, pulled)); err != nil {
					logger.Error().Err(err).Msg("error while handling config")
				}
				nextDuration := time.Until(pulled.NextPullAfter())
				if nextDuration < 5*time.Minute {
					nextDuration = 5 * time.Minute
				}
				nextDuration = jitter(nextDuration)
				logger.Debug().Dur("nextPoolInMillis", nextDuration).Msg("next config poll")
				timer.Reset(nextDuration)
			}
		}
	}
	if err := p.cache.Persist(false); err != nil {
		logger.Error().Err(err).Msg("error while persisting config")
	}
	return nil
}

func (p *ConfigPoller[CK, T]) ReceivePushedConfig(c context.Context, config T) (context.Context, error) {
	return p.Handle(context.WithValue(c, p.contextKey, config))
}

func (p *ConfigPoller[CK, T]) Name() string {
	return p.name
}

func NewConfigPoller[CK comparable, T VersionedConfig](
	name string, contextKey CK,
	pullRemoteConfig func(context.Context) (T, error),
	cache *ConfigCache[T]) *ConfigPoller[CK, T] {
	return &ConfigPoller[CK, T]{
		name:             name,
		contextKey:       contextKey,
		pullRemoteConfig: pullRemoteConfig,
		cache:            cache,
	}
}

var _ taskmanager.Task = (*ConfigPoller[string, VersionedConfig])(nil)
