package common

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/shared"
)

var BuildID string = "dev"

func RespondDiagnostics(c echo.Context) error {
	r := shared.RequestDiagnostics{
		ServiceRuntime: shared.ServiceRuntimeInfo{
			GoVersion: runtime.Version(),
			BuildID:   BuildID,
		},
	}
	for k, v := range c.Request().Header {
		r.RequestHeaders = append(r.RequestHeaders, shared.RequestHeaderEntry{
			Key:   k,
			Value: v,
		})
	}
	return c.JSON(http.StatusOK, r)
}
