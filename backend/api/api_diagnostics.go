package api

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/models"
)

func (*server) GetDiagnostics(c echo.Context) error {
	r := models.RequestDiagnostics{
		ServiceRuntime: models.RequestDiagnostics_ServiceRuntime{
			GoVersion: runtime.Version(),
		},
	}
	for k, v := range c.Request().Header {
		r.RequestHeaders = append(r.RequestHeaders, models.RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	return c.JSON(http.StatusOK, r)
}
