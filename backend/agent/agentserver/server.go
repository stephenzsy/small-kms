package agentserver

import (
	"net/http"

	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type server struct {
	common.CommonServer
	dockerClient *dockerclient.Client
}

// GetDockerInfo implements ServerInterface.
func (s *server) GetDockerInfo(ctx echo.Context, _ shared.Identifier) error {
	info, err := s.dockerClient.Info(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, info)
}

func NewServer() (*server, error) {

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		return nil, err
	}

	config, err := common.NewCommonConfig()
	if err != nil {
		return nil, err
	}

	s := &server{
		CommonServer: config,
		dockerClient: cli,
	}

	return s, err
}
