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

func (s *proxiedServer) delegateRequest(ec echo.Context,
	namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier,
	resourceIdentifier base.Identifier, delegatedAuthToken *string,
	getResult func(ctx.RequestContext, ClientWithResponsesInterface) ([]byte, int, error)) error {
	c := ec.(ctx.RequestContext)

	if delegatedAuthToken == nil || *delegatedAuthToken == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, resourceIdentifier, *delegatedAuthToken)
	if err != nil {
		return err
	}
	body, status, err := getResult(c, client)
	if err != nil {
		return err
	}
	return c.JSONBlob(status, body)
}

// AgentDockerContainerInspect implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerInspect(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, containerId string, params AgentDockerContainerInspectParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, resourceIdentifier, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) ([]byte, int, error) {
		result, err := client.AgentDockerContainerInspectWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, containerId, nil)
		if err != nil {
			return nil, 0, err
		}
		return result.Body, result.StatusCode(), nil
	})
}

// AgentDockerContainerList implements ServerInterface.
func (s *proxiedServer) AgentDockerContainerList(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params AgentDockerContainerListParams) error {
	return s.delegateRequest(ec, namespaceKind, namespaceIdentifier, resourceIdentifier, params.XCryptocatProxyAuthorization, func(c ctx.RequestContext, client ClientWithResponsesInterface) ([]byte, int, error) {
		result, err := client.AgentDockerContainerListWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, nil)
		if err != nil {
			return nil, 0, err
		}
		return result.Body, result.StatusCode(), nil
	})
}

// AgentDockerImageList implements ServerInterface.
func (s *proxiedServer) AgentDockerImageList(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params AgentDockerImageListParams) error {
	c := ec.(ctx.RequestContext)

	if params.XCryptocatProxyAuthorization == nil || *params.XCryptocatProxyAuthorization == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, resourceIdentifier, *params.XCryptocatProxyAuthorization)
	if err != nil {
		return err
	}
	resp, err := client.AgentDockerImageListWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, nil)
	if err != nil {
		return err
	}
	return c.JSONBlob(resp.StatusCode(), resp.Body)
}

// AgentPullImage implements ServerInterface.
func (s *proxiedServer) AgentPullImage(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params AgentPullImageParams) error {
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
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, resourceIdentifier, *params.XCryptocatProxyAuthorization)
	if err != nil {
		return err
	}
	resp, err := client.AgentPullImageWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, nil, req)
	if err != nil {
		return err
	}
	return c.JSONBlob(resp.StatusCode(), resp.Body)
}

// GetDiagnostics implements ServerInterface.
func (s *proxiedServer) GetAgentDiagnostics(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params GetAgentDiagnosticsParams) error {
	c := ec.(ctx.RequestContext)

	if params.XCryptocatProxyAuthorization == nil || *params.XCryptocatProxyAuthorization == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, resourceIdentifier, *params.XCryptocatProxyAuthorization)
	if err != nil {
		return err
	}
	resp, err := client.GetAgentDiagnosticsWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, nil)
	if err != nil {
		return err
	}
	return c.JSONBlob(resp.StatusCode(), resp.Body)
}

// GetDockerInfo implements ServerInterface.
func (s *proxiedServer) GetAgentDockerInfo(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier, params GetAgentDockerInfoParams) error {
	c := ec.(ctx.RequestContext)

	if params.XCryptocatProxyAuthorization == nil || *params.XCryptocatProxyAuthorization == "" {
		return fmt.Errorf("%w: missing delegated access token", base.ErrResposneStatusUnauthorized)
	}

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	// proxiedClient
	client, err := s.getProxiedClient(c, resourceIdentifier, *params.XCryptocatProxyAuthorization)
	if err != nil {
		return err
	}
	resp, err := client.GetAgentDockerInfoWithResponse(c, namespaceKind, namespaceIdentifier, resourceIdentifier, nil)
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
