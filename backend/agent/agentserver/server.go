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
	}
	apiBaseUrl, ok := os.LookupEnv("SMALLKMS_API_BASE_URL")
	if !ok {

		return nil, fmt.Errorf("environment variable SMALLKMS_API_BASE_URL is not set")
	}
	apiScope, ok := os.LookupEnv("SMALLKMS_API_SCOPE")
	if !ok {
		return nil, fmt.Errorf("environment variable SMALLKMS_API_SCOPE is not set")
	}
	tenantID, ok := os.LookupEnv("AZURE_TENANT_ID")
	if !ok {
		return nil, fmt.Errorf("environment variable AZURE_TENANT_ID is not set")
	}
	configDir, ok := os.LookupEnv("SMALLKMS_AGENT_CONFIG_DIR")
	if !ok {
		return nil, fmt.Errorf("environment variable SMALLKMS_AGENT_CONFIG_DIR is not set")
	}
	azKeyVaultUrl, ok := os.LookupEnv("AZURE_KEYVAULT_RESOURCEENDPOINT")
	if !ok {
		return nil, fmt.Errorf("environment variable SMALLKMS_AGENT_CONFIG_DIR is not set")
	}
	s.ConfigLoader, err = newConfigLoader(
		config.ServiceIdentity(),
		apiBaseUrl,
		apiScope,
		tenantID,
		configDir,
		azKeyVaultUrl,
	)

	return s, err
}
