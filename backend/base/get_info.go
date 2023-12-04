package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
)

func RespondDiagnostics(c ctx.RequestContext, info models.ServiceRuntimeInfo) error {
	r := models.RequestDiagnostics{
		ServiceRuntime:  info,
		RequestProtocol: c.Request().Proto,
	}
	for k, v := range c.Request().Header {
		r.RequestHeaders = append(r.RequestHeaders, models.RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	return c.JSON(http.StatusOK, r)
}

func RespondHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
