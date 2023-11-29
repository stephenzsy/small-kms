package agentpush

import (
	"fmt"

	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type proxiedServer struct {
	api.APIServer
	proxyClientPool *ProxyClientPool
}

// PushAgentConfigRadius implements ServerInterface.
func (s *proxiedServer) PushAgentConfigRadius(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, rID base.ID, params PushAgentConfigRadiusParams) error {
	return s.delegateRequest(ec, nsKind, nsID, rID, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
		return client.PushAgentConfigRadiusWithResponse(c, nsKind, nsID, rID, nil)
	})
}

// AgentContainerRemove implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerRemove(ec echo.Context, namespaceKind base.NamespaceKind,
	nsID base.ID,
	rID base.ID,
	containerId string, params AgentDockerContainerRemoveParams) error {
	return s.delegateRequest(ec, namespaceKind, nsID, rID, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
		return client.AgentDockerContainerRemoveWithResponse(c, namespaceKind, nsID, rID, containerId, nil)
	})
}

// AgentDockerContainerStop implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerStop(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	containerId string, params AgentDockerContainerStopParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, rID, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
		return client.AgentDockerContainerStopWithResponse(c, namespaceKind, namespaceIdentifier, rID, containerId, nil)
	})
}

// AgentLaunchAgent implements ServerInterface.
func (s *proxiedServer) AgentLaunchAgent(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	params AgentLaunchAgentParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, rID, params.XCryptocatProxyAuthorization,
		func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
			return client.AgentLaunchAgentWithBodyWithResponse(
				c, namespaceKind,
				namespaceIdentifier,
				rID,
				nil,
				c.Request().Header.Get(echo.HeaderContentType),
				c.Request().Body)
		})
}

func (s *proxiedServer) delegateRequest(ec echo.Context,
	namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	delegatedAuthToken *string,
	getResult func(ctx.RequestContext, ClientWithResponsesInterface) (ProxiedResponse, error)) error {
	c := ec.(ctx.RequestContext)

	if delegatedAuthToken == nil || *delegatedAuthToken == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)
	client, err := s.getProxiedClient(c, rID, *delegatedAuthToken)
	if err != nil {
		return err
	}
	result, err := getResult(c, client)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(result.StatusCode(), result.GetBody())
}

// AgentDockerContainerInspect implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerInspect(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	containerId string, params AgentDockerContainerInspectParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, rID, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
		return client.AgentDockerContainerInspectWithResponse(c, namespaceKind, namespaceIdentifier, rID, containerId, nil)
	})
}

// AgentDockerContainerList implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerList(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	params AgentDockerContainerListParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, rID, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) (ProxiedResponse, error) {
		return client.AgentDockerContainerListWithResponse(c, namespaceKind, namespaceIdentifier, rID, nil)
	})
}

// AgentPullImage implements ServerInterface.
func (s *proxiedServer) AgentPullImage(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID,
	rID base.ID,
	params AgentPullImageParams) error {
	c := ec.(ctx.RequestContext)

	if params.XCryptocatProxyAuthorization == nil || *params.XCryptocatProxyAuthorization == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}
	req := PullImageRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, rID, *params.XCryptocatProxyAuthorization)
	if err != nil {
		return err
	}
	resp, err := client.AgentPullImageWithResponse(c, namespaceKind, namespaceIdentifier, rID, nil, req)
	if err != nil {
		return err
	}
	return c.JSONBlob(resp.StatusCode(), resp.Body)
}

func NewProxiedServer(apiServer api.APIServer) ServerInterface {
	return &proxiedServer{
		APIServer:       apiServer,
		proxyClientPool: NewProxyClientPool(128),
	}
}
