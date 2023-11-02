package agentpush

import (
	"net/http"
	"runtime"

	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
)

type agentServer struct {
	common.CommonServer
	buildID      string
	dockerClient *dockerclient.Client
}

// GetDiagnostics implements ServerInterface.
func (s *agentServer) GetAgentDiagnostics(c echo.Context, _ base.NamespaceKind, _, _ base.Identifier, _ GetAgentDiagnosticsParams) error {
	return base.RespondDiagnostics(c, base.ServiceRuntimeInfo{
		BuildID:   s.buildID,
		GoVersion: runtime.Version(),
	})
}

// GetDockerInfo implements ServerInterface.
func (s *agentServer) GetAgentDockerInfo(ctx echo.Context, _ base.NamespaceKind, _, _ base.Identifier, _ GetAgentDockerInfoParams) error {

	info, err := s.dockerClient.Info(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, info)
}

var _ ServerInterface = (*agentServer)(nil)

func NewServer(buildID string) (*agentServer, error) {

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		return nil, err
	}

	config, err := common.NewCommonConfig()
	if err != nil {
		return nil, err
	}

	s := &agentServer{
		CommonServer: config,
		buildID:      buildID,
		dockerClient: cli,
	}

	return s, err
}
