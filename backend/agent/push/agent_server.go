package agentpush

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type agentServer struct {
	common.CommonServer
	buildID         string
	dockerClient    *dockerclient.Client
	acrAuthProvider *acr.DockerRegistryAuthProvider
	acrImageRepo    string
}

// AgentDockerImageList implements ServerInterface.
func (s *agentServer) AgentDockerImageList(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params AgentDockerImageListParams) error {
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
func (s *agentServer) AgentPullImage(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params AgentPullImageParams) error {
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

	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetDiagnostics implements ServerInterface.
func (s *agentServer) GetAgentDiagnostics(ec echo.Context, _ base.NamespaceKind, _, _ base.Identifier, _ GetAgentDiagnosticsParams) error {
	c := ec.(ctx.RequestContext)

	return base.RespondDiagnostics(c, base.ServiceRuntimeInfo{
		BuildID:   s.buildID,
		GoVersion: runtime.Version(),
	})
}

// GetDockerInfo implements ServerInterface.
func (s *agentServer) GetAgentDockerInfo(ec echo.Context, _ base.NamespaceKind, _, _ base.Identifier, _ GetAgentDockerInfoParams) error {
	c := ec.(ctx.RequestContext)

	info, err := s.dockerClient.Info(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

var _ ServerInterface = (*agentServer)(nil)

func NewServer(buildID string) (*agentServer, error) {
	var acrLoginServer string
	var tenantID string
	var acrImageRepo string
	var err error
	if tenantID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "AZURE_TENANT_ID", ""); tenantID == "" {
		return nil, fmt.Errorf("%w:%s", common.ErrInvalidEnvVar, "AZURE_TENANT_ID")
	} else if acrImageRepo = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "AZURE_ACR_IMAGE_REPOSITORY", ""); acrImageRepo == "" {
		return nil, fmt.Errorf("%w:%s", common.ErrInvalidEnvVar, "AZURE_ACR_IMAGE_REPOSITORY")
	} else if acrLoginServer, err = acr.ExtractACRLoginServer(acrImageRepo); err != nil {
		return nil, err
	}

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		return nil, err
	}

	config, err := common.NewCommonConfig()
	if err != nil {
		return nil, err
	}

	s := &agentServer{
		CommonServer:    config,
		buildID:         buildID,
		dockerClient:    cli,
		acrAuthProvider: acr.NewDockerRegistryAuthProvider(acrLoginServer, config.ServiceIdentity().TokenCredential(), tenantID),
		acrImageRepo:    acrImageRepo,
	}

	return s, err
}