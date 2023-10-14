package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *server) GetDiagnostics(c echo.Context) error {
	return common.RespondDiagnostics(c, s.getRuntimeInfo())
}
