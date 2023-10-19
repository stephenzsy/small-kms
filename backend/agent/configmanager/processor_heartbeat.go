package cm

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type heartbeatConfigProcessor struct {
	baseConfigProcessor
}

// Version implements ConfigProcessor.
func (*heartbeatConfigProcessor) Version() string {
	return "1.0"
}

func (p *heartbeatConfigProcessor) Process(ctx context.Context, _ string) error {
	resp, err := p.AgentClient().AgentCallbackWithResponse(ctx,
		shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
		shared.AgentConfigNameHeartbeat,
		shared.AgentConfiguration{})
	if err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msgf("heartbeat response: %d", resp.StatusCode())
	return nil
}

func (p *heartbeatConfigProcessor) Start(c context.Context, scheduleToUpdate chan<- pollConfigMsg, shutdownNotifier common.LeafShutdownNotifier) {
	p.baseStart(c, scheduleToUpdate, func() *pollConfigMsg {
		return &pollConfigMsg{
			name:      shared.AgentConfigNameHeartbeat,
			processor: p,
		}
	}, func() {}, shutdownNotifier)
}

var _ ConfigProcessor = (*heartbeatConfigProcessor)(nil)

func (p *heartbeatConfigProcessor) MarkProcessDone(string, error) {
	resetTimer := p.baseMarkProcessDone()
	resetTimer(5*time.Minute, "heartbeat")
}

func newHeartbeatConfigProcessor(sc *sharedConfig) *heartbeatConfigProcessor {
	return &heartbeatConfigProcessor{
		baseConfigProcessor: baseConfigProcessor{
			sharedConfig: sc,
			configName:   shared.AgentConfigNameHeartbeat,
		},
	}
}
