package configmanager

import "context"

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
	if c, err = h.Before(c); err != nil {
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
