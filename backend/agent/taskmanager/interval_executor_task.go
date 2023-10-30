package taskmanager

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type IntervalExecutor interface {
	Name() string
	Execute(context.Context) (time.Duration, error)
}

type intervalExecutorTask struct {
	executor     IntervalExecutor
	initialDelay time.Duration
}

// Name implements Task.
func (t *intervalExecutorTask) Name() string {
	return t.executor.Name()
}

// Start implements Task.
func (t *intervalExecutorTask) Start(c context.Context, sigCh <-chan os.Signal) error {
	logger := log.Ctx(c).With().Str("task", t.executor.Name()).Logger()
	logger.Debug().Msg("Task started")
	defer logger.Debug().Msg("Task exited in interval executor")

	timer := time.NewTimer(t.initialDelay)
	defer timer.Stop()
	active := true
	var err error

	for active {
		select {
		case <-c.Done():
			return c.Err()
		case <-sigCh:
			active = false
		case <-timer.C:
			var nextDelay time.Duration
			nextDelay, err = t.executor.Execute(c)
			logger.Debug().Dur("next execute", nextDelay).Err(err).Msgf("Executed")
			if err != nil {
				logger.Error().Err(err).Msgf("Task %s returned error", t.executor.Name())
			}
			if active {
				timer.Reset(nextDelay)
			}
		}
	}
	// return last error
	return err
}

func IntervalExecutorTask(executor IntervalExecutor, initialDelay time.Duration) Task {
	return &intervalExecutorTask{
		executor: executor,
	}
}
