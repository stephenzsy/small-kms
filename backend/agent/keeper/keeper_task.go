package keeper

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/agent/taskmanager"
)

type keeperTaskExecutor struct {
	cm *ConfigManager
}

// Execute implements taskmanager.ScheduledTask.
func (t *keeperTaskExecutor) Execute(c context.Context) (time.Duration, error) {
	if !t.cm.IsReady() {
		err := t.cm.PullConfig(c)
		log.Ctx(c).Debug().Err(err).Msg("Pull config")
	}
	return time.Minute * 5, nil
}

func (*keeperTaskExecutor) Name() string {
	return "Keeper"
}

func NewKeeper(configManager *ConfigManager) taskmanager.IntervalExecutor {
	return &keeperTaskExecutor{
		cm: configManager,
	}
}
