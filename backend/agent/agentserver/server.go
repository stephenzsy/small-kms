package agentserver

import (
	"net/http"

	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
)

type server struct {
	dockerClient *dockerclient.Client
}

// GetDockerInfo implements ServerInterface.
func (s *server) GetDockerInfo(ctx echo.Context) error {
	info, err := s.dockerClient.Info(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, info)
}

func NewServer() ServerInterface {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		panic(err)
	}

	return &server{
		dockerClient: cli,
	}
}
