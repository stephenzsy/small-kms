package base

import (
	"net/http"
	"runtime"

	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type baseServer struct {
	common.CommonServer
}

// GetDiagnostics implements ServerInterface.
func (s *baseServer) GetDiagnostics(ec echo.Context) error {
	c := ec.(ctx.RequestContext)
	return RespondDiagnostics(c, s.getRuntimeInfo())
}

// GetHealth implements ServerInterface.
func (*baseServer) GetHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func (s *baseServer) getRuntimeInfo() ServiceRuntimeInfo {
	return ServiceRuntimeInfo{
		BuildID:     s.BuildID(),
		GoVersion:   runtime.Version(),
		Environment: s.EnvService().Export(),
	}
}

// NewBaseServer creates a new base server.
func NewBaseServer(server common.CommonServer) ServerInterface {
	return &baseServer{
		CommonServer: server,
	}
}
