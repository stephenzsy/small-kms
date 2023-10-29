package base

import (
	"net/http"
	"runtime"

	echo "github.com/labstack/echo/v4"
)

type baseServer struct {
	buildID string
}

// GetDiagnostics implements ServerInterface.
func (s *baseServer) GetDiagnostics(ctx echo.Context) error {
	return RespondDiagnostics(ctx, s.getRuntimeInfo())
}

// GetHealth implements ServerInterface.
func (*baseServer) GetHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func (s *baseServer) getRuntimeInfo() ServiceRuntimeInfo {
	return ServiceRuntimeInfo{
		BuildID:   s.buildID,
		GoVersion: runtime.Version(),
	}
}

// NewBaseServer creates a new base server.
func NewBaseServer(buildID string) ServerInterface {
	return &baseServer{
		buildID: buildID,
	}
}
