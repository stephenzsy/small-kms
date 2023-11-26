package base

import (
	"net/http"
	"runtime"

	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
)

type BaseServer struct {
	common.CommonServer
}

// GetDiagnostics implements ServerInterface.
func (s *BaseServer) GetDiagnostics(ec echo.Context) error {
	c := ec.(ctx.RequestContext)
	return RespondDiagnostics(c, s.getRuntimeInfo())
}

// GetHealth implements ServerInterface.
func (*BaseServer) GetHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func (s *BaseServer) getRuntimeInfo() models.ServiceRuntimeInfo {
	return models.ServiceRuntimeInfo{
		BuildID:     s.BuildID(),
		GoVersion:   runtime.Version(),
		Environment: s.EnvService().Export(),
	}
}

// NewBaseServer creates a new base server.
func NewBaseServer(server common.CommonServer) *BaseServer {
	return &BaseServer{
		CommonServer: server,
	}
}
