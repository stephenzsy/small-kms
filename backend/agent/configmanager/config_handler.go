package configmanager

import "context"

type ContextConfigHandler interface {
	Before(context.Context) (context.Context, error)
	After(context.Context) (context.Context, error)
	Handle(context.Context) (context.Context, error)
}

type ChainedContextConfigHandler struct {
	next *ChainedContextConfigHandler
}

// After implements ContextConfigHandler.
func (*ChainedContextConfigHandler) After(c context.Context) (context.Context, error) {
	return c, nil
}

// Before implements ContextConfigHandler.
func (*ChainedContextConfigHandler) Before(c context.Context) (context.Context, error) {
	return c, nil
}

// Before implements ContextConfigHandler.
func (h *ChainedContextConfigHandler) Handle(c context.Context) (context.Context, error) {
	if h == nil {
		return c, nil
	}
	c, err := h.Before(c)
	if err != nil {
		return c, err
	}
	if h.next != nil {
		c, err := h.next.Handle(c)
		if err != nil {
			return c, err
		}
	}
	return h.After(c)
}

var _ ContextConfigHandler = (*ChainedContextConfigHandler)(nil)
