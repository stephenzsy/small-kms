package common

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type shutdownNotifier struct {
	quit     chan os.Signal
	complete chan struct{}
}

func NewStandaloneShutdownNotifier() *shutdownNotifier {
	return &shutdownNotifier{
		quit:     make(chan os.Signal, 1),
		complete: make(chan struct{}, 1),
	}
}

func NewLeafShutdownNotifier() LeafShutdownNotifier {
	return NewStandaloneShutdownNotifier()
}

func (n *shutdownNotifier) Quit() <-chan os.Signal {
	return n.quit
}

func (n *shutdownNotifier) Complete() <-chan struct{} {
	return n.complete
}

func (n *shutdownNotifier) ListenOSIntercept() {
	signal.Notify(n.quit, os.Interrupt)
}

func (n *shutdownNotifier) RelaySingal(sig os.Signal) {
	n.quit <- sig
}

func (n *shutdownNotifier) MarkShutdownComplete() {
	n.complete <- struct{}{}
	close(n.complete)
}
func (n *shutdownNotifier) RelayComplete(m MergedShutdownNotifier) {
	n.complete <- <-m.Complete()
	close(n.complete)
}

type MergedShutdownNotifier interface {
	Complete() <-chan struct{}
	ListenOSIntercept()
	RelaySingal(sig os.Signal)
}

type LeafShutdownNotifier interface {
	Quit() <-chan os.Signal
	MarkShutdownComplete()
	RelayComplete(MergedShutdownNotifier)
}

func MergeShutdownNotifier(notifiers ...LeafShutdownNotifier) MergedShutdownNotifier {
	if len(notifiers) == 0 {
		return nil
	}
	n := NewStandaloneShutdownNotifier()
	go func() {
		for sig := range n.quit {
			for _, v := range notifiers {
				v.(*shutdownNotifier).quit <- sig
			}
		}
		for _, v := range notifiers {
			close(v.(*shutdownNotifier).quit)
		}
	}()
	go func() {
		for _, v := range notifiers {
			<-v.(*shutdownNotifier).complete
		}
		n.complete <- struct{}{}
		close(n.complete)
	}()
	return n
}

func StartEchoWithGracefulShutdown(c context.Context, e *echo.Echo, onStart func(*echo.Echo, LeafShutdownNotifier), gracePeriod time.Duration) {
	c = log.Logger.WithContext(c)
	shutdownNotifier := NewStandaloneShutdownNotifier()

	go onStart(e, shutdownNotifier)

	shutdownNotifier.ListenOSIntercept()
	<-shutdownNotifier.Quit()
	toCtx, toCancel := context.WithTimeout(c, gracePeriod)
	defer toCancel()
	go func() {
		log.Info().Msg("echo server shutdown")
		err := e.Shutdown(toCtx)
		log.Info().Err(err).Msg("echo server shutted down")
	}()

	select {
	case <-toCtx.Done():
		log.Fatal().Err(toCtx.Err()).Msg("failed to gracefully shutdown")
	case <-shutdownNotifier.Quit():
		log.Fatal().Msg("forced shutdown after second interupt")
	case <-shutdownNotifier.Complete():
		log.Info().Msg("gracefully shutdown")
	}
}
