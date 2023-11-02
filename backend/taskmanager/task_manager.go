package taskmanager

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type Task interface {
	Name() string
	Start(c context.Context, sigCh <-chan os.Signal) error
}

type taskImpl struct {
	name  string
	start func(c context.Context, sigCh <-chan os.Signal) error
}

func (t *taskImpl) Name() string {
	return t.name
}

func (t *taskImpl) Start(c context.Context, sigCh <-chan os.Signal) error {
	return t.start(c, sigCh)
}

func NewTask(name string, start func(c context.Context, sigCh <-chan os.Signal) error) Task {
	return &taskImpl{
		name:  name,
		start: start,
	}
}

type TaskManager interface {
	Start(c context.Context, sigCh <-chan os.Signal) error
	WithTask(Task) TaskManager
}

type chainedTaskManager struct {
	next *chainedTaskManager
	task Task
}

// WithTask implements TaskManager.
func (ctm *chainedTaskManager) WithTask(t Task) TaskManager {
	if ctm.task == nil {
		ctm.task = t
	} else if ctm.next != nil {
		ctm.next.WithTask(t)
	} else {
		// ctm.next == nil
		ctm.next = &chainedTaskManager{task: t}
	}
	return ctm
}

type taskController struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	sigCh  chan os.Signal
	done   chan error
}

func newTaskController(c context.Context) *taskController {
	innerCtx, innerCtxCancel := context.WithCancelCause(c)
	return &taskController{
		ctx:    innerCtx,
		cancel: innerCtxCancel,
		sigCh:  make(chan os.Signal, 1),
		done:   make(chan error, 1),
	}
}

// Start implements TaskManager.
func (ctm *chainedTaskManager) Start(ctx context.Context, sigCh <-chan os.Signal) error {
	logger := log.Ctx(ctx).With().Str("task", ctm.task.Name()).Logger()
	logger.Debug().Msg("Starting task")
	defer logger.Debug().Msg("Task exited")

	var taskCtrl, nextCtrl *taskController
	if ctm.task != nil {
		taskCtrl = newTaskController(ctx)
		go func() {
			defer close(taskCtrl.done)
			taskCtrl.done <- ctm.task.Start(taskCtrl.ctx, taskCtrl.sigCh)
		}()
	}
	if ctm.next != nil {
		nextCtrl = newTaskController(ctx)
		go func() {
			defer close(nextCtrl.done)
			nextCtrl.done <- ctm.next.Start(nextCtrl.ctx, nextCtrl.sigCh)
		}()
	}

	active := true
	var err error
	for active {
		select {
		case <-ctx.Done():
			active = false
			if taskCtrl != nil {
				taskCtrl.cancel(ctx.Err())
			}
			if nextCtrl != nil {
				nextCtrl.cancel(ctx.Err())
			}
			return ctx.Err()
		case sig, sigOpen := <-sigCh:
			active = false
			if taskCtrl != nil {
				taskCtrl.sigCh <- sig
				if !sigOpen {
					close(taskCtrl.sigCh)
				}
			}
			if nextCtrl != nil {
				nextCtrl.sigCh <- sig
				if !sigOpen {
					close(nextCtrl.sigCh)
				}
			}
			if nextCtrl != nil {
				err = <-nextCtrl.done
			}
			if taskCtrl != nil {
				taskErr := <-taskCtrl.done
				if err != nil {
					err = fmt.Errorf("task %s exited with errors: %w, %s", ctm.task.Name(), err, taskErr)
				}
				logger.Debug().Err(taskErr).Msgf("Task %s exited", ctm.task.Name())
			}
		}
	}
	return err
}

var _ TaskManager = (*chainedTaskManager)(nil)

func NewChainedTaskManager() TaskManager {
	return &chainedTaskManager{}
}
