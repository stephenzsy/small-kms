package cm

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type baseConfigProcessor struct {
	*sharedConfig
	configName     shared.AgentConfigName
	isActive       bool
	pollShutdownCh chan struct{}
	processPending bool
	timer          *time.Timer
	attemptedLoad  bool
}

// Name implements ConfigProcessor.
func (p *baseConfigProcessor) ConfigName() shared.AgentConfigName {
	return p.configName
}

func (p *baseConfigProcessor) baseStart(c context.Context, scheduleToUpdate chan<- pollConfigMsg,
	exitCh chan<- error,
	onTimer func() *pollConfigMsg) func() {
	p.timer = time.NewTimer(0)
	p.isActive = true
	for p.isActive {
		select {
		case <-p.pollShutdownCh:
			p.isActive = false
		case <-c.Done():
			exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrForcedShutdown, p.ConfigName(), c.Err())
			return func() {}
		case <-p.timer.C:
			if !p.processPending && p.isActive {
				p.processPending = true
				msg := onTimer()
				if msg != nil {
					scheduleToUpdate <- *msg
				}
			}
		}
	}
	return func() {
		exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrGracefullyShutdown, p.ConfigName(), c.Err())
	}
}

func (p *baseConfigProcessor) baseShutdown() func() {
	p.isActive = false
	p.timer.Stop()
	return func() {
		p.pollShutdownCh <- struct{}{}
	}
}

func (p *baseConfigProcessor) Shutdown() {
	cb := p.baseShutdown()
	cb()
}

func (p *baseConfigProcessor) baseMarkProcessDone() func(time.Duration) {
	p.processPending = true
	p.timer.Stop()
	return func(d time.Duration) {
		log.Info().Msgf("%s: next refresh after: %s", p.configName, d)
		p.processPending = false
		p.timer.Reset(d)
	}
}
