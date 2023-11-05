package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func RespondDiagnostics(c ctx.RequestContext, info ServiceRuntimeInfo) error {
	r := RequestDiagnostics{
		ServiceRuntime: info,
	}
	for k, v := range c.Request().Header {
		r.RequestHeaders = append(r.RequestHeaders, RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	return c.JSON(http.StatusOK, r)
}

func RespondHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
