package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/models"
)

// AgentCheckIn implements models.ServerInterface.
func (*server) AgentCheckIn(ctx echo.Context, params models.AgentCheckInParams) error {
	return ctx.NoContent(http.StatusNoContent)
}
