package agentadmin

import (
	"github.com/labstack/echo/v4"
	agentendpoint "github.com/stephenzsy/small-kms/backend/agent/endpoint"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type AgentPushProxiedServer struct {
	api.APIServer
	clientPool *ProxyClientPool
}

// AgentDockerImagePull implements agentendpoint.ServerInterface.
func (s *AgentPushProxiedServer) AgentDockerImagePull(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	req := new(agentmodels.PullImageRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	client, err := s.clientPool.GetClient(c, namespaceId, id)
	if err != nil {
		return err
	}
	resp, err := client.AgentDockerImagePullWithResponse(c, "me", "me", *req)
	if err != nil {
		return err
	}
	return c.NoContent(resp.StatusCode())
}

// ListAgentDockerNetowks implements agentendpoint.ServerInterface.
func (s *AgentPushProxiedServer) ListAgentDockerNetowks(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	client, err := s.clientPool.GetClient(c, namespaceId, id)
	if err != nil {
		return err
	}
	resp, err := client.ListAgentDockerNetowksWithResponse(c, "me", "me")
	if err != nil {
		return err
	}
	return c.Blob(resp.StatusCode(), "application/json", resp.Body)
}

// GetAgentDiagnostics implements admin.ServerInterface.
func (s *AgentPushProxiedServer) GetAgentDiagnostics(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	client, err := s.clientPool.GetClient(c, namespaceId, id)
	if err != nil {
		return err
	}
	resp, err := client.GetAgentDiagnosticsWithResponse(c, "me", "me")
	if err != nil {
		return err
	}
	return c.Blob(resp.StatusCode(), "application/json", resp.Body)
}

// GetAgentDockerSystemInformation implements admin.ServerInterface.
func (s *AgentPushProxiedServer) GetAgentDockerSystemInformation(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	client, err := s.clientPool.GetClient(c, namespaceId, id)
	if err != nil {
		return err
	}
	resp, err := client.GetAgentDockerSystemInformationWithResponse(c, "me", "me")
	if err != nil {
		return err
	}
	return c.Blob(resp.StatusCode(), "application/json", resp.Body)
}

// ListAgentDockerImages implements agentendpoint.ServerInterface.
func (s *AgentPushProxiedServer) AgentDockerImageList(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	client, err := s.clientPool.GetClient(c, namespaceId, id)
	if err != nil {
		return err
	}
	resp, err := client.AgentDockerImageListWithResponse(c, "me", "me")
	if err != nil {
		return err
	}
	return c.Blob(resp.StatusCode(), "application/json", resp.Body)
}

var _ agentendpoint.ServerInterface = (*AgentPushProxiedServer)(nil)

func NewAgentPushProxiedServer(apiServer api.APIServer) *AgentPushProxiedServer {
	return &AgentPushProxiedServer{
		APIServer:  apiServer,
		clientPool: NewProxyClientPool(16),
	}
}
