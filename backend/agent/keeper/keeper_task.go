package keeper

import (
	"context"
	"time"

	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type keeperTaskExecutor struct {
	cm           *ConfigManager
	configUpdate chan AgentServerConfiguration
}

// BeforeShutdown implements taskmanager.IntervalExecutor.
func (e *keeperTaskExecutor) Close(c context.Context) error {
	close(e.configUpdate)
	e.cm.configProcessor.Shutdown(c)
	return nil
}

// Execute implements taskmanager.ScheduledTask.
func (t *keeperTaskExecutor) Execute(c context.Context) (time.Duration, error) {
	bad := func(err error) (time.Duration, error) {
		return time.Minute * 5, err
	}
	if !t.cm.HasAttemptedLoad() {
		config, err := t.cm.LoadConfig(c)
		if err != nil {
			return bad(err)
		} else if config == nil {
			return 5 * time.Second, nil
		}
		t.configUpdate <- config
		return config.NextWaitInterval(), nil
	}
	config, hasChange, err := t.cm.PullConfig(c)
	if err != nil {
		return bad(err)
	} else if hasChange {
		t.configUpdate <- config
	}
	return config.NextWaitInterval(), nil
}

func (*keeperTaskExecutor) Name() string {
	return "Keeper"
}

var _ taskmanager.IntervalExecutor = (*keeperTaskExecutor)(nil)

func NewKeeper(configManager *ConfigManager) *keeperTaskExecutor {
	return &keeperTaskExecutor{
		cm:           configManager,
		configUpdate: make(chan AgentServerConfiguration),
	}
}

func (e *keeperTaskExecutor) ConfigUpdate() <-chan AgentServerConfiguration {
	return e.configUpdate
}
