package agentconfigmanager

import (
	"context"
	"time"

	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type pollTaskExecutor struct {
	cm *configManager
}

// BeforeShutdown implements taskmanager.IntervalExecutor.
func (e *pollTaskExecutor) Close(c context.Context) error {
	return nil
}

// Execute implements taskmanager.ScheduledTask.
func (t *pollTaskExecutor) Execute(c context.Context) (time.Duration, error) {
	nextPullAt, err := t.cm.pullConfig(c)
	nextInterval := time.Until(nextPullAt)
	if nextInterval < time.Minute*5 {
		nextInterval = time.Minute * 5
	}
	return nextInterval, err
}

func (*pollTaskExecutor) Name() string {
	return "ConfigManager2"
}

var _ taskmanager.IntervalExecutor = (*pollTaskExecutor)(nil)

func NewConfigManagerPollTaskExecutor(cm *configManager) *pollTaskExecutor {
	return &pollTaskExecutor{
		cm: cm,
	}
}
