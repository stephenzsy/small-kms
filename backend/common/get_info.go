package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var BuildID = "dev"

func RespondDiagnostics(c echo.Context, info shared.ServiceRuntimeInfo) error {
	r := shared.RequestDiagnostics{
		ServiceRuntime: info,
	}
	for k, v := range c.Request().Header {
		r.RequestHeaders = append(r.RequestHeaders, shared.RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	return c.JSON(http.StatusOK, r)
}
