package agentpush

import (
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/keeper"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

type agentServer struct {
	common.CommonServer
	buildID             string
	dockerClient        *dockerclient.Client
	acrAuthProvider     *acr.DockerRegistryAuthProvider
	acrImageRepo        string
	mode                managedapp.AgentMode
	radiusConfigManager *keeper.RadiusConfigManager
}

// AgentContainerRemove implements ServerInterface.
func (s *agentServer) AgentDockerContainerRemove(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, containerId string, params AgentDockerContainerRemoveParams) error {
	c := ec.(ctx.RequestContext)

	err := s.dockerClient.ContainerRemove(c, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// AgentDockerContainerStop implements ServerInterface.
func (s *agentServer) AgentDockerContainerStop(ec echo.Context, _ base.NamespaceKind, _, _ base.ID,
	containerId string, params AgentDockerContainerStopParams) error {
	c := ec.(ctx.RequestContext)
	return s.apiStopContainer(c, containerId)
}

// AgentLaunchContainer implements ServerInterface.
func (s *agentServer) AgentLaunchAgent(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentLaunchAgentParams) error {
	c := ec.(ctx.RequestContext)
	req := LaunchAgentRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	return s.apiLaunchAgentContainer(c, req)
}

// AgentDockerContainerInspect implements ServerInterface.
func (s *agentServer) AgentDockerContainerInspect(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, containerId string, _ AgentDockerContainerInspectParams) error {
	c := ec.(ctx.RequestContext)

	result, err := s.dockerClient.ContainerInspect(c, containerId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// AgentDockerImageList implements ServerInterface.
func (s *agentServer) AgentDockerContainerList(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentDockerContainerListParams) error {
	c := ec.(ctx.RequestContext)

	containers, err := s.dockerClient.ContainerList(c, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, containers)
}

// AgentDockerImageList implements ServerInterface.
func (s *agentServer) AgentDockerImageList(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentDockerImageListParams) error {
	c := ec.(ctx.RequestContext)

	images, err := s.dockerClient.ImageList(c, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, images)
}

// AgentPullImage implements ServerInterface.
func (s *agentServer) AgentPullImage(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentPullImageParams) error {
	c := ec.(ctx.RequestContext)
	req := PullImageRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	regAuth, err := s.acrAuthProvider.GetRegistryAuth(c)
	if err != nil {
		return err
	}
	if req.ImageTag == "" {
		return fmt.Errorf("%w: missing image tag", base.ErrResponseStatusBadRequest)
	}
	imageRef := fmt.Sprintf("%s:%s", s.acrImageRepo, req.ImageTag)
	out, err := s.dockerClient.ImagePull(c, imageRef, types.ImagePullOptions{
		RegistryAuth: regAuth,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(io.Discard, out)
	return c.NoContent(http.StatusNoContent)
}

func (s *agentServer) AgentDockerNetworkList(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentDockerNetworkListParams) error {
	c := ec.(ctx.RequestContext)

	result, err := s.dockerClient.NetworkList(c, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// GetDiagnostics implements ServerInterface.
func (s *agentServer) GetAgentDiagnostics(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ GetAgentDiagnosticsParams) error {
	c := ec.(ctx.RequestContext)

	return base.RespondDiagnostics(c, base.ServiceRuntimeInfo{
		BuildID:     s.buildID,
		GoVersion:   runtime.Version(),
		Environment: s.EnvService().Export(),
	})
}

// GetDockerInfo implements ServerInterface.
func (s *agentServer) AgentDockerInfo(ec echo.Context, _ base.NamespaceKind, _, _ base.ID, _ AgentDockerInfoParams) error {
	c := ec.(ctx.RequestContext)

	info, err := s.dockerClient.Info(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

var _ ServerInterface = (*agentServer)(nil)

func NewServer(buildID string, mode managedapp.AgentMode, envSvc common.EnvService, radiusConfigManager *keeper.RadiusConfigManager) (*agentServer, error) {
	var acrLoginServer string
	var tenantID string
	var acrImageRepo string
	var err error
	var ok bool
	if tenantID, ok = envSvc.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, envSvc.ErrMissing(common.EnvKeyAzTenantID)
	} else if acrImageRepo, ok = envSvc.RequireNonWhitespace(agentcommon.EnvKeyAcrImageRepository, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, envSvc.ErrMissing(agentcommon.EnvKeyAcrImageRepository)
	} else if acrLoginServer, err = acr.ExtractACRLoginServer(acrImageRepo); err != nil {
		return nil, err
	}

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		return nil, err
	}

	config, err := common.NewCommonConfig(envSvc, buildID)
	if err != nil {
		return nil, err
	}

	s := &agentServer{
		CommonServer:        config,
		buildID:             buildID,
		dockerClient:        cli,
		acrAuthProvider:     acr.NewDockerRegistryAuthProvider(acrLoginServer, config.ServiceIdentity().TokenCredential(), tenantID),
		acrImageRepo:        acrImageRepo,
		mode:                mode,
		radiusConfigManager: radiusConfigManager,
	}

	return s, err
}
