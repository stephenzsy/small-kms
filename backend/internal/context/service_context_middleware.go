package ctx

import (
	"context"

	"github.com/labstack/echo/v4"
)

func InjectServiceContextMiddleware(serviceContext context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c = NewInjectedRequestContext(c, serviceContext)
			return next(c)
		}
	}
}
