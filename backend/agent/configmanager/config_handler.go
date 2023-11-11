package configmanager

import (
	"context"

	"github.com/rs/zerolog/log"
)

type HandleContextConfigFunc func(context.Context) (context.Context, error)

type ContextConfigHandler interface {
	Before(c context.Context) (context.Context, error)
	After(context.Context) (context.Context, error)
}

type ChainedContextConfigHandler struct {
	ContextConfigHandler
	next *ChainedContextConfigHandler
}

func (h *ChainedContextConfigHandler) Handle(c context.Context) (context.Context, error) {
	var err error
	logger := log.Ctx(c)
	if c, err = h.Before(c); err != nil {
		logger.Error().Err(err).Msg("failed to handle config")
		return c, err
	}
	if h.next != nil {
		if c, err = h.next.Handle(c); err != nil {
			return c, err
		}
	}
	return h.After(c)
}

func (h *ChainedContextConfigHandler) SetNext(nh *ChainedContextConfigHandler) {
	nh.next = h.next
	h.next = nh
}

var _ ContextConfigHandler = (*ChainedContextConfigHandler)(nil)
