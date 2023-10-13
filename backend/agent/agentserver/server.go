package agentserver

import (
	"fmt"
	"net/http"
	"os"

	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/common"
)

type server struct {
	common.CommonServer
	dockerClient *dockerclient.Client
	ConfigLoader
}

// GetDockerInfo implements ServerInterface.
func (s *server) GetDockerInfo(ctx echo.Context) error {
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
		ConfigLoader: ConfigLoader{
			identity:      config.ServiceIdentity(),
			cacheFileName: "agent-config.json",
		},
	}
	if apiBaseUrl, ok := os.LookupEnv("SMALLKMS_API_BASE_URL"); ok {
		s.ConfigLoader.baseUrl = apiBaseUrl
	} else {
		return nil, fmt.Errorf("environment variable SMALLKMS_API_BASE_URL is not set")
	}
	if apiScope, ok := os.LookupEnv("SMALLKMS_API_SCOPE"); ok {
		s.ConfigLoader.authScope = apiScope
	} else {
		return nil, fmt.Errorf("environment variable SMALLKMS_API_SCOPE is not set")
	}
	if tenantID, ok := os.LookupEnv("AZURE_TENANT_ID"); ok {
		s.ConfigLoader.tenantID = tenantID
	} else {
		return nil, fmt.Errorf("environment variable AZURE_TENANT_ID is not set")
	}
	return s, nil
}
