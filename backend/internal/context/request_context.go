package ctx

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type internalRequestContextKey int

const (
	rKey internalRequestContextKey = iota
)

type (
	RequestContext interface {
		echo.Context
		context.Context
		WithValue(key any, val any) RequestContext
	}

	wrappedCtx struct {
		context.Context
	}

	requestContext struct {
		echo.Context
		elevated   bool
		reqCtx     *wrappedCtx
		serviceCtx *wrappedCtx
	}
)

func (c *requestContext) chCtx() *wrappedCtx {
	if c.elevated {
		return c.serviceCtx
	}
	return c.reqCtx
}

// Deadline implements RequestContext.
func (c *requestContext) Deadline() (time.Time, bool) {
	return c.chCtx().Deadline()
}

// Done implements RequestContext.
func (c *requestContext) Done() <-chan struct{} {
	return c.chCtx().Done()
}

// Err implements RequestContext.
func (c *requestContext) Err() error {
	return c.chCtx().Err()
}

// Value implements RequestContext.
func (c *requestContext) Value(key any) any {
	switch key {
	case rKey:
		return c
	}
	if val := c.reqCtx.Value(key); val != nil {
		return val
	}
	return c.serviceCtx.Value(key)
}

var _ RequestContext = (*requestContext)(nil)

func (c *requestContext) WithValue(key any, val any) RequestContext {
	return &requestContext{
		Context:    c.Context,
		elevated:   c.elevated,
		reqCtx:     &wrappedCtx{context.WithValue(c.reqCtx.Context, key, val)},
		serviceCtx: c.serviceCtx,
	}
}

func (c *requestContext) withServiceContexValue(key any, val any) *requestContext {
	return &requestContext{
		Context:    c.Context,
		elevated:   c.elevated,
		reqCtx:     c.reqCtx,
		serviceCtx: &wrappedCtx{context.WithValue(c.serviceCtx.Context, key, val)},
	}
}

func NewInjectedRequestContext(c echo.Context, serviceCtx context.Context) RequestContext {
	reqCtx := &wrappedCtx{Context: c.Request().Context()}
	return &requestContext{
		Context:    c,
		elevated:   false,
		reqCtx:     reqCtx,
		serviceCtx: &wrappedCtx{serviceCtx},
	}
}

func Elevate(c context.Context) context.Context {
	if reqCtx, ok := c.Value(rKey).(*requestContext); ok && reqCtx != nil {
		if reqCtx.elevated {
			// already elevated
			return c
		}
		nextReqCtx := reqCtx.reqCtx
		if reqCtx != c {
			nextReqCtx = &wrappedCtx{c}
		}
		return &requestContext{
			Context:    reqCtx.Context,
			elevated:   true,
			reqCtx:     nextReqCtx,
			serviceCtx: reqCtx.serviceCtx,
		}
	}
	return c
}

func EchoContextWithValue(c echo.Context, key any, value any, isServiceContext bool) echo.Context {
	if reqCtx, ok := c.(*requestContext); ok && reqCtx != nil {
		if isServiceContext {
			return reqCtx.withServiceContexValue(key, value)
		} else {
			return reqCtx.WithValue(key, value)
		}
	}
	return c
}
