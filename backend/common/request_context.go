package common

import (
	ctx "context"
	"time"

	"github.com/labstack/echo/v4"
)

type ElevatedRequestContext interface {
	ServerContext() ServerContext
}

type RequestContext interface {
	echo.Context
	ctx.Context
	IsElevated() bool
	Elevate() RequestContext
	ServiceClientProvider() ServiceClientProvider
}

type injectedEchoContext struct {
	echo.Context
	inner      ctx.Context
	isElevated bool
}

// ServiceClientProvider implements RequestContext.
func (c injectedEchoContext) ServiceClientProvider() ServiceClientProvider {
	return c.Value(serverContextKey).(ServerContext)
}

// IsElevated implements RequestContext.
func (c injectedEchoContext) IsElevated() bool {
	return c.isElevated
}

// Elevate implements RequestContext.
func (c injectedEchoContext) Elevate() RequestContext {
	return injectedEchoContext{
		Context:    c.Context,
		inner:      c.Value(serverContextKey).(ServerContext),
		isElevated: true,
	}
}

// Deadline implements InjectedEchoContext.
func (c injectedEchoContext) Deadline() (deadline time.Time, ok bool) {
	return c.inner.Deadline()
}

// Done implements InjectedEchoContext.
func (c injectedEchoContext) Done() <-chan struct{} {
	return c.inner.Done()
}

// Err implements InjectedEchoContext.
func (c injectedEchoContext) Err() error {
	return c.inner.Err()
}

// Value implements InjectedEchoContext.
func (c injectedEchoContext) Value(key any) any {
	return c.inner.Value(key)
}

var _ RequestContext = injectedEchoContext{}

func WrapEchoContext(c echo.Context) RequestContext {
	return injectedEchoContext{
		Context: c,
		inner:   c.Request().Context(),
	}
}

func EchoContextWithValue(parent echo.Context, key any, val any) RequestContext {
	if p, ok := parent.(RequestContext); ok {
		return RequestContextWithValue(p, key, val)
	}
	return injectedEchoContext{
		Context: parent,
		inner:   ctx.WithValue(parent.Request().Context(), key, val),
	}
}

type contextKey string

const serverContextKey contextKey = "serverContext"

func EchoContextWithServerContext(parent echo.Context, val ServerContext) RequestContext {
	return EchoContextWithValue(parent, serverContextKey, val)
}

func RequestContextWithValue(parent RequestContext, key any, val any) RequestContext {
	return injectedEchoContext{
		Context: parent,
		inner:   ctx.WithValue(parent, key, val),
	}
}

func InjectAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(WrapEchoContext(c))
	}
}
