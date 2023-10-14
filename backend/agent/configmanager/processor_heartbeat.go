package cm

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type heartbeatConfigProcessor struct {
	*sharedConfig
	isActive       bool
	timer          *time.Timer
	processPending bool
	pollShutdownCh chan struct{}
}

// Name implements ConfigProcessor.
func (*heartbeatConfigProcessor) Name() shared.AgentConfigName {
	return shared.AgentConfigNameHeartbeat
}

var meNamespaceIdIdentifier = shared.StringIdentifier("me")

func (p *heartbeatConfigProcessor) Process(ctx context.Context) error {
	resp, err := p.AgentClient().AgentCallbackWithResponse(ctx,
		shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
		shared.AgentConfigNameHeartbeat,
		shared.AgentCallbackRequest{
			ServiceRuntime: p.sharedConfig.ServiceRuntime(),
		})
	if err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msgf("heartbeat response: %d", resp.StatusCode())
	return nil
}

func (p *heartbeatConfigProcessor) Start(c context.Context, scheduleToUpdate chan<- pollConfigMsg, exitCh chan<- error) {
	p.isActive = true
	p.timer = time.NewTimer(0)
	defer p.timer.Stop()
	for p.isActive {
		select {
		case <-p.pollShutdownCh:
			goto pollEnd
		case <-c.Done():
			exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrForcedShutdown, p.Name(), c.Err())
			return
		case <-p.timer.C:
			if !p.processPending && p.isActive {
				p.processPending = true
				scheduleToUpdate <- pollConfigMsg{
					name:      shared.AgentConfigNameHeartbeat,
					processor: p,
				}
			}
		}
	}
pollEnd:
	exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrGracefullyShutdown, p.Name(), c.Err())
}

var _ ConfigProcessor = (*heartbeatConfigProcessor)(nil)

func (p *heartbeatConfigProcessor) Shutdown() {
	p.isActive = false
	p.timer.Stop()
	p.pollShutdownCh <- struct{}{}
}

func (p *heartbeatConfigProcessor) MarkProcessDone() {
	p.timer.Stop()
	p.timer.Reset(5 * time.Minute)
	p.processPending = false
}

func newHeartbeatConfigProcessor(sc *sharedConfig) *heartbeatConfigProcessor {
	return &heartbeatConfigProcessor{
		sharedConfig:   sc,
		pollShutdownCh: make(chan struct{}, 1),
	}
}
