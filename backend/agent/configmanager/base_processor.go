package cm

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type baseConfigProcessor struct {
	*sharedConfig
	configName     shared.AgentConfigName
	processPending bool
	timer          *time.Timer
	attemptedLoad  bool
	shutdownCtrl   *ShutdownController
}

// Name implements ConfigProcessor.
func (p *baseConfigProcessor) ConfigName() shared.AgentConfigName {
	return p.configName
}

func (p *baseConfigProcessor) baseStart(c context.Context, scheduleToUpdate chan<- pollConfigMsg,
	onTimer func() *pollConfigMsg,
	beforeShutdown func()) {
	shutdownCh, shutdownDefer := p.shutdownCtrl.Subsribe()
	defer shutdownDefer()
	p.timer = time.NewTimer(0)
	for p.shutdownCtrl.IsActive() {
		select {
		case <-shutdownCh:
		case <-c.Done():
			log.Error().Err(c.Err()).Msgf("processor:%s forced shutdown", p.ConfigName())
			return
		case <-p.timer.C:
			if !p.processPending && p.shutdownCtrl.IsActive() {
				p.processPending = true
				msg := onTimer()
				if msg != nil {
					scheduleToUpdate <- *msg
				}
			}
		}
	}
	beforeShutdown()
}

func (p *baseConfigProcessor) baseMarkProcessDone() func(time.Duration, string) {
	p.processPending = true
	p.timer.Stop()
	return func(d time.Duration, taskName string) {
		log.Info().Msgf("%s:%s next refresh after: %s", p.configName, taskName, d)
		p.processPending = false
		p.timer.Reset(d)
	}
}
