package agentpush

import (
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	agentendpoint "github.com/stephenzsy/small-kms/backend/agent/endpoint"
	"github.com/stephenzsy/small-kms/backend/agent/radius"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type agentServer struct {
	*base.BaseServer
	dockerClient        dockerclient.APIClient
	acrAuthProvider     *acr.DockerRegistryAuthProvider
	acrImageRepo        string
	mode                agentcommon.AgentSlot
	radiusConfigManager *radius.RadiusConfigManager
	launchedBy          string
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
func (s *agentServer) ListAgentDockerImages(ec echo.Context, namespaceId string, id string) error {
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
	imageRepo := s.acrImageRepo
	if req.ImageRepo != "" {
		loginServer, err := acr.ExtractACRLoginServer(req.ImageRepo)
		if err != nil {
			return err
		}
		reqLoginServer, err := acr.ExtractACRLoginServer(req.ImageRepo)
		if err != nil {
			return err
		}
		if loginServer != reqLoginServer {
			return fmt.Errorf("%w: image repo must be in the same as configured", base.ErrResponseStatusBadRequest)
		}
		imageRepo = req.ImageRepo
	}

	imageRef := fmt.Sprintf("%s:%s", imageRepo, req.ImageTag)
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

func (s *agentServer) ListAgentDockerNetowks(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	result, err := s.dockerClient.NetworkList(c, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// GetAgentDiagnostics implements agentendpoint.ServerInterface.
func (s *agentServer) GetAgentDiagnostics(ec echo.Context, _ string, _ string) error {
	return s.BaseServer.GetDiagnostics(ec)
}

// GetDockerInfo implements ServerInterface.
func (s *agentServer) GetAgentDockerSystemInformation(ec echo.Context, _, _ string) error {
	c := ec.(ctx.RequestContext)

	info, err := s.dockerClient.Info(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

var _ ServerInterface = (*agentServer)(nil)
var _ agentendpoint.ServerInterface = (*agentServer)(nil)

func NewServer(buildID string, mode agentcommon.AgentSlot, envSvc common.EnvService, dockerClient dockerclient.APIClient) (*agentServer, error) {
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

	config, err := common.NewCommonConfig(envSvc, buildID)
	if err != nil {
		return nil, err
	}

	s := &agentServer{
		BaseServer:      base.NewBaseServer(config),
		dockerClient:    dockerClient,
		acrAuthProvider: acr.NewDockerRegistryAuthProvider(acrLoginServer, config.ServiceIdentity().TokenCredential(), tenantID),
		acrImageRepo:    acrImageRepo,
		mode:            mode,
		launchedBy:      envSvc.Default("AGENT_LAUNCHED_BY", ""),
	}

	return s, err
}
