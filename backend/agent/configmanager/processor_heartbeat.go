package cm

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type heartbeatConfigProcessor struct {
	baseConfigProcessor
	timer *time.Timer
}

// Version implements ConfigProcessor.
func (*heartbeatConfigProcessor) Version() string {
	return "1.0"
}

func (p *heartbeatConfigProcessor) Process(ctx context.Context, _ string) error {
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
	p.timer = time.NewTimer(0)
	p.baseStart(c, scheduleToUpdate, exitCh, p.timer.C, func() *pollConfigMsg {
		return &pollConfigMsg{
			name:      shared.AgentConfigNameHeartbeat,
			processor: p,
		}
	})
}

var _ ConfigProcessor = (*heartbeatConfigProcessor)(nil)

func (p *heartbeatConfigProcessor) Shutdown() {
	baseDefer := p.baseShutdown()
	defer baseDefer()
	p.timer.Stop()
}

func (p *heartbeatConfigProcessor) MarkProcessDone(string, error) {
	baseDefer := p.baseMarkProcessDone()
	defer baseDefer()
	p.timer.Stop()
	p.timer.Reset(5 * time.Minute)
}

func newHeartbeatConfigProcessor(sc *sharedConfig) *heartbeatConfigProcessor {
	return &heartbeatConfigProcessor{
		baseConfigProcessor: baseConfigProcessor{
			sharedConfig:   sc,
			configName:     shared.AgentConfigNameHeartbeat,
			pollShutdownCh: make(chan struct{}, 1),
		},
	}
}
