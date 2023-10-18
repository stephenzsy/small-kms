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
	name     string
	quit     chan os.Signal
	complete chan struct{}
}

func NewStandaloneShutdownNotifier(name string) *shutdownNotifier {
	return &shutdownNotifier{
		name:     name,
		quit:     make(chan os.Signal, 1),
		complete: make(chan struct{}, 1),
	}
}

func NewLeafShutdownNotifier(name string) LeafShutdownNotifier {
	return NewStandaloneShutdownNotifier(name)
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
	//log.Debug().Msgf("mark shutdown complete: %s", n.name)
}

func (n *shutdownNotifier) RelayComplete(m MergedShutdownNotifier) {
	n.complete <- <-m.Complete()
	close(n.complete)
	//log.Debug().Msgf("shutdown relay complete: %s", n.name)
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

func MergeShutdownNotifier(name string, notifiers ...LeafShutdownNotifier) MergedShutdownNotifier {
	if len(notifiers) == 0 {
		return nil
	}
	n := NewStandaloneShutdownNotifier(name)
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
			//log.Debug().Msgf("%s: waiting shutdown complete: %s", n.name, v.(*shutdownNotifier).name)
			<-v.(*shutdownNotifier).complete
			//log.Debug().Msgf("%s: received shutdown complete: %s", n.name, v.(*shutdownNotifier).name)
		}
		n.MarkShutdownComplete()
	}()
	return n
}

func StartEchoWithGracefulShutdown(c context.Context, e *echo.Echo, onStart func(*echo.Echo, LeafShutdownNotifier), gracePeriod time.Duration) {
	c = log.Logger.WithContext(c)
	shutdownNotifier := NewStandaloneShutdownNotifier("echo server")

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
