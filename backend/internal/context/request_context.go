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
	}

	wrappedCtx struct {
		context.Context
	}

	requestContext struct {
		echo.Context
		chCtx      *wrappedCtx
		valueCtx   *wrappedCtx
		serviceCtx *wrappedCtx
	}
)

// Deadline implements RequestContext.
func (c *requestContext) Deadline() (time.Time, bool) {
	return c.chCtx.Deadline()
}

// Done implements RequestContext.
func (c *requestContext) Done() <-chan struct{} {
	return c.chCtx.Done()
}

// Err implements RequestContext.
func (c *requestContext) Err() error {
	return c.chCtx.Err()
}

// Value implements RequestContext.
func (c *requestContext) Value(key any) any {
	switch key {
	case rKey:
		return c
	}
	if val := c.valueCtx.Value(key); val != nil {
		return val
	}
	return c.serviceCtx.Value(key)
}

var _ RequestContext = (*requestContext)(nil)

func (c *requestContext) withValue(key any, val any) *requestContext {
	if key == rKey {
		return c
	}
	nextValueCtx := &wrappedCtx{context.WithValue(c.valueCtx.Context, key, val)}
	nextChCtx := c.chCtx
	if nextChCtx == c.valueCtx {
		nextChCtx = nextValueCtx
	}
	return &requestContext{
		Context:    c.Context,
		chCtx:      nextChCtx,
		valueCtx:   nextValueCtx,
		serviceCtx: c.serviceCtx,
	}
}

func NewInjectedRequestContext(c echo.Context, serviceCtx context.Context) RequestContext {
	reqCtx := &wrappedCtx{Context: c.Request().Context()}
	return &requestContext{
		Context:    c,
		chCtx:      reqCtx,
		valueCtx:   reqCtx,
		serviceCtx: &wrappedCtx{serviceCtx},
	}
}

func Elevate(c context.Context) context.Context {
	if reqCtx, ok := c.Value(rKey).(*requestContext); ok && reqCtx != nil {
		if reqCtx.chCtx == reqCtx.serviceCtx {
			// already elevated
			return c
		}
		valueCtx := reqCtx.valueCtx
		if reqCtx != c {
			valueCtx = &wrappedCtx{c}
		}
		return &requestContext{
			Context:    reqCtx.Context,
			chCtx:      reqCtx.serviceCtx,
			valueCtx:   valueCtx,
			serviceCtx: reqCtx.serviceCtx,
		}
	}
	return c
}

func EchoContextWithValue(c echo.Context, key any, value any) echo.Context {
	if reqCtx, ok := c.(*requestContext); ok && reqCtx != nil {
		return reqCtx.withValue(key, value)
	}
	return c
}

func ResolveRequestContext(c echo.Context) RequestContext {
	if cc, ok := c.(context.Context); ok && cc != nil {
		return cc.Value(rKey).(RequestContext)
	}
	return nil
}
