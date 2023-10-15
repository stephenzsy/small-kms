package cm

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type activeServerProcessor struct {
	baseConfigProcessor
	timer *time.Timer
}

// Version implements ConfigProcessor.
func (*activeServerProcessor) Version() string {
	return "0.1"
}

// Name implements ConfigProcessor.
func (*activeServerProcessor) Name() shared.AgentConfigName {
	return shared.AgentConfigNameHeartbeat
}

func (p *activeServerProcessor) Process(ctx context.Context, task string) error {
	switch task {
	case TaskNameLoad:
		if err := p.loadFetchedConfig(); err != nil {
			return err
		}
	case TaskNameFetch:
		resp, err := p.AgentClient().GetAgentConfigurationWithResponse(ctx,
			shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
			shared.AgentConfigNameActiveServer,
			nil)
		if err != nil {
			return err
		}
		nextConfig := resp.JSON200
		log.Ctx(ctx).Info().Msgf("fetched config: %s:%s", p.configName, nextConfig.Version)
		if p.fetchedConfig == nil || nextConfig.Version != p.fetchedConfig.Version {
			if err := p.saveNewFetchedVersion(nextConfig); err != nil {
				// this error is not fatal
				log.Ctx(ctx).Error().Err(err).Msgf("save new fetched version: %s", nextConfig.Version)
			}
			p.fetchedConfig = nextConfig
		}
	default:
		return fmt.Errorf("unknown task: %s", task)
	}
	return nil
}

func (p *activeServerProcessor) Start(c context.Context, scheduleToUpdate chan<- pollConfigMsg, exitCh chan<- error) {
	p.timer = time.NewTimer(0)
	attemptedLoading := false
	p.baseStart(c, scheduleToUpdate, exitCh, p.timer.C, func() *pollConfigMsg {
		if !attemptedLoading {
			attemptedLoading = true
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameLoad,
				processor: p,
			}
		}
		if p.fetchedConfig == nil || p.fetchedConfig.NextRefreshAfter == nil || time.Now().After(*p.fetchedConfig.NextRefreshAfter) {
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameFetch,
				processor: p,
			}
		}
		return nil
	})
}

var _ ConfigProcessor = (*activeServerProcessor)(nil)

func (p *activeServerProcessor) Shutdown() {
	baseDefer := p.baseShutdown()
	defer baseDefer()
	p.timer.Stop()
}

func (p *activeServerProcessor) MarkProcessDone(taskName string, err error) {
	baseDefer := p.baseMarkProcessDone()
	defer baseDefer()
	nextDuration := 5 * time.Minute
	switch taskName {
	case TaskNameLoad:
		p.timer.Stop()
		if err != nil || p.fetchedConfig == nil || p.fetchedConfig.NextRefreshAfter == nil {
			// load failed, wait 5 seconds before fetching from server
			nextDuration = 5 * time.Second
		} else {
			nextDuration = time.Until(*p.fetchedConfig.NextRefreshAfter) + 59*time.Second
		}
	case TaskNameFetch:
		p.timer.Stop()
		if err != nil || p.fetchedConfig == nil || p.fetchedConfig.NextRefreshAfter == nil {
			// fetch failed, wait 5 minutes before retrying
		} else {
			nextDuration = time.Until(*p.fetchedConfig.NextRefreshAfter) + 59*time.Second
		}
	}
	if nextDuration < 0 {
		// if next duration is less than 0, wait 1 hour
		nextDuration = time.Hour
	}
	log.Info().Msgf("%s:next refresh after: %s", p.configName, nextDuration)
	p.timer.Reset(nextDuration)
}

func newActiveServerProcessor(sc *sharedConfig) *activeServerProcessor {
	return &activeServerProcessor{
		baseConfigProcessor: baseConfigProcessor{
			sharedConfig:   sc,
			configName:     shared.AgentConfigNameActiveServer,
			pollShutdownCh: make(chan struct{}, 1),
		},
	}
}
