package taskmanager

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

func StartWithGracefulShutdown(c context.Context, tm TaskManager) error {
	logger := log.Ctx(c)
	sigRec := make(chan os.Signal, 2)
	signal.Notify(sigRec, os.Interrupt)

	var exitErr error
	ctxWithCancel, cancel := context.WithCancelCause(c)
	defer cancel(exitErr)

	sigSend := make(chan os.Signal, 1)
	defer close(sigSend)

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- tm.Start(ctxWithCancel, sigSend)
	}()

	sigReceived := false
	active := true

	for active {
		select {
		case <-c.Done():
			active = false
			logger.Debug().Msg("Context cancelled, forece exit")
			exitErr = c.Err()
		case sig := <-sigRec:
			if sigReceived {
				active = false
				// second signal received, force exit
				exitErr = errors.New("second signal received, force exit")
				logger.Debug().Msg("Second signal received, force exit")
				cancel(exitErr)
			} else {
				sigReceived = true
				logger.Debug().Msgf("Received signal: %s", sig)
				sigSend <- sig
				var shutdownTimeoutCancel context.CancelFunc
				c, shutdownTimeoutCancel = context.WithTimeout(c, time.Second*20)
				defer shutdownTimeoutCancel()
			}
		case exitErr = <-doneCh:
			active = false
		}
	}
	return exitErr
}
