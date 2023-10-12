package common

import (
	ctx "context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RequestContext struct {
	echo.Context
	valueCtx       ctx.Context
	serviceContext ctx.Context
}

// Deadline implements context.Context.
func (c RequestContext) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Request().Context().Deadline()
}

// Done implements context.Context.
func (c RequestContext) Done() <-chan struct{} {
	return c.Context.Request().Context().Done()
}

// Err implements context.Context.
func (c RequestContext) Err() error {
	return c.Context.Request().Context().Err()
}

// Value implements context.Context.
func (c RequestContext) Value(key any) any {
	if v := c.valueCtx.Value(key); v != nil {
		return v
	}
	if v := c.serviceContext.Value(key); v != nil {
		return v
	}
	return c.Context.Request().Context().Value(key)
}

func (c RequestContext) Elevate() ctx.Context {
	return ctx.WithValue(c.serviceContext, isElevatedContextKey, true)
}

func IsElevated(c ctx.Context) bool {
	if v := c.Value(isElevatedContextKey); v != nil {
		return v.(bool)
	}
	return false
}

func (c RequestContext) WithSharedValue(key any, val any) RequestContext {
	return RequestContext{
		Context:        c.Context,
		valueCtx:       c.valueCtx,
		serviceContext: ctx.WithValue(c.serviceContext, key, val),
	}
}

func (c RequestContext) WitValue(key any, val any) RequestContext {
	return RequestContext{
		Context:        c.Context,
		valueCtx:       ctx.WithValue(c.valueCtx, key, val),
		serviceContext: c.serviceContext,
	}
}

var _ echo.Context = RequestContext{}
var _ ctx.Context = RequestContext{}

func WrapEchoContext(c echo.Context, serviceContext ctx.Context) RequestContext {
	return RequestContext{
		Context:        c,
		valueCtx:       ctx.Background(),
		serviceContext: serviceContext,
	}
}

type H = map[string]string

func EchoRequestContextWithSharedValue(c echo.Context, key any, val any) (echo.Context, error) {
	if r, ok := c.(RequestContext); ok {
		return r.WithSharedValue(key, val), nil
	}
	log.Error().Msg("invalid echo context in chain")
	return c, c.JSON(http.StatusInternalServerError, H{"error": "internal error"})
}
