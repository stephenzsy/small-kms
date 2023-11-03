// Package agentpush provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package agentpush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAgentDiagnostics request
	GetAgentDiagnostics(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDiagnosticsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AgentDockerContainerList request
	AgentDockerContainerList(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerContainerListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AgentDockerContainerInspect request
	AgentDockerContainerInspect(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, containerId string, params *AgentDockerContainerInspectParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AgentDockerImageList request
	AgentDockerImageList(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerImageListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAgentDockerInfo request
	GetAgentDockerInfo(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDockerInfoParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AgentPullImageWithBody request with any body
	AgentPullImageWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AgentPullImage(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, body AgentPullImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAgentDiagnostics(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDiagnosticsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAgentDiagnosticsRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentDockerContainerList(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerContainerListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentDockerContainerListRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentDockerContainerInspect(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, containerId string, params *AgentDockerContainerInspectParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentDockerContainerInspectRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, containerId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentDockerImageList(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerImageListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentDockerImageListRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAgentDockerInfo(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDockerInfoParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAgentDockerInfoRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentPullImageWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentPullImageRequestWithBody(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentPullImage(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, body AgentPullImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentPullImageRequest(c.Server, namespaceKind, namespaceIdentifier, resourceIdentifier, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAgentDiagnosticsRequest generates requests for GetAgentDiagnostics
func NewGetAgentDiagnosticsRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDiagnosticsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/diagnostics", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

// NewAgentDockerContainerListRequest generates requests for AgentDockerContainerList
func NewAgentDockerContainerListRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerContainerListParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/docker/containers", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

// NewAgentDockerContainerInspectRequest generates requests for AgentDockerContainerInspect
func NewAgentDockerContainerInspectRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, containerId string, params *AgentDockerContainerInspectParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam3 string

	pathParam3, err = runtime.StyleParamWithLocation("simple", false, "containerId", runtime.ParamLocationPath, containerId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/docker/containers/%s", pathParam0, pathParam1, pathParam2, pathParam3)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

// NewAgentDockerImageListRequest generates requests for AgentDockerImageList
func NewAgentDockerImageListRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerImageListParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/docker/images", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

// NewGetAgentDockerInfoRequest generates requests for GetAgentDockerInfo
func NewGetAgentDockerInfoRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDockerInfoParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/docker/info", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

// NewAgentPullImageRequest calls the generic AgentPullImage builder with application/json body
func NewAgentPullImageRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, body AgentPullImageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAgentPullImageRequestWithBody(server, namespaceKind, namespaceIdentifier, resourceIdentifier, params, "application/json", bodyReader)
}

// NewAgentPullImageRequestWithBody generates requests for AgentPullImage with any type of body
func NewAgentPullImageRequestWithBody(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, namespaceIdentifier)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, resourceIdentifier)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s/pull-image", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.XCryptocatProxyAuthorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Cryptocat-Proxy-Authorization", runtime.ParamLocationHeader, *params.XCryptocatProxyAuthorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Cryptocat-Proxy-Authorization", headerParam0)
		}

	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAgentDiagnosticsWithResponse request
	GetAgentDiagnosticsWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDiagnosticsParams, reqEditors ...RequestEditorFn) (*GetAgentDiagnosticsResponse, error)

	// AgentDockerContainerListWithResponse request
	AgentDockerContainerListWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerContainerListParams, reqEditors ...RequestEditorFn) (*AgentDockerContainerListResponse, error)

	// AgentDockerContainerInspectWithResponse request
	AgentDockerContainerInspectWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, containerId string, params *AgentDockerContainerInspectParams, reqEditors ...RequestEditorFn) (*AgentDockerContainerInspectResponse, error)

	// AgentDockerImageListWithResponse request
	AgentDockerImageListWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerImageListParams, reqEditors ...RequestEditorFn) (*AgentDockerImageListResponse, error)

	// GetAgentDockerInfoWithResponse request
	GetAgentDockerInfoWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDockerInfoParams, reqEditors ...RequestEditorFn) (*GetAgentDockerInfoResponse, error)

	// AgentPullImageWithBodyWithResponse request with any body
	AgentPullImageWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AgentPullImageResponse, error)

	AgentPullImageWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, body AgentPullImageJSONRequestBody, reqEditors ...RequestEditorFn) (*AgentPullImageResponse, error)
}

type GetAgentDiagnosticsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.RequestDiagnostics
}

// Status returns HTTPResponse.Status
func (r GetAgentDiagnosticsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAgentDiagnosticsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AgentDockerContainerListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]DockerContainer
}

// Status returns HTTPResponse.Status
func (r AgentDockerContainerListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AgentDockerContainerListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AgentDockerContainerInspectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]DockerContainerJSON
}

// Status returns HTTPResponse.Status
func (r AgentDockerContainerInspectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AgentDockerContainerInspectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AgentDockerImageListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]DockerImageSummary
}

// Status returns HTTPResponse.Status
func (r AgentDockerImageListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AgentDockerImageListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAgentDockerInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *DockerInfo
}

// Status returns HTTPResponse.Status
func (r GetAgentDockerInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAgentDockerInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AgentPullImageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r AgentPullImageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AgentPullImageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAgentDiagnosticsWithResponse request returning *GetAgentDiagnosticsResponse
func (c *ClientWithResponses) GetAgentDiagnosticsWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDiagnosticsParams, reqEditors ...RequestEditorFn) (*GetAgentDiagnosticsResponse, error) {
	rsp, err := c.GetAgentDiagnostics(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAgentDiagnosticsResponse(rsp)
}

// AgentDockerContainerListWithResponse request returning *AgentDockerContainerListResponse
func (c *ClientWithResponses) AgentDockerContainerListWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerContainerListParams, reqEditors ...RequestEditorFn) (*AgentDockerContainerListResponse, error) {
	rsp, err := c.AgentDockerContainerList(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentDockerContainerListResponse(rsp)
}

// AgentDockerContainerInspectWithResponse request returning *AgentDockerContainerInspectResponse
func (c *ClientWithResponses) AgentDockerContainerInspectWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, containerId string, params *AgentDockerContainerInspectParams, reqEditors ...RequestEditorFn) (*AgentDockerContainerInspectResponse, error) {
	rsp, err := c.AgentDockerContainerInspect(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, containerId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentDockerContainerInspectResponse(rsp)
}

// AgentDockerImageListWithResponse request returning *AgentDockerImageListResponse
func (c *ClientWithResponses) AgentDockerImageListWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentDockerImageListParams, reqEditors ...RequestEditorFn) (*AgentDockerImageListResponse, error) {
	rsp, err := c.AgentDockerImageList(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentDockerImageListResponse(rsp)
}

// GetAgentDockerInfoWithResponse request returning *GetAgentDockerInfoResponse
func (c *ClientWithResponses) GetAgentDockerInfoWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *GetAgentDockerInfoParams, reqEditors ...RequestEditorFn) (*GetAgentDockerInfoResponse, error) {
	rsp, err := c.GetAgentDockerInfo(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAgentDockerInfoResponse(rsp)
}

// AgentPullImageWithBodyWithResponse request with arbitrary body returning *AgentPullImageResponse
func (c *ClientWithResponses) AgentPullImageWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AgentPullImageResponse, error) {
	rsp, err := c.AgentPullImageWithBody(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentPullImageResponse(rsp)
}

func (c *ClientWithResponses) AgentPullImageWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter, params *AgentPullImageParams, body AgentPullImageJSONRequestBody, reqEditors ...RequestEditorFn) (*AgentPullImageResponse, error) {
	rsp, err := c.AgentPullImage(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentPullImageResponse(rsp)
}

// ParseGetAgentDiagnosticsResponse parses an HTTP response from a GetAgentDiagnosticsWithResponse call
func ParseGetAgentDiagnosticsResponse(rsp *http.Response) (*GetAgentDiagnosticsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAgentDiagnosticsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest externalRef0.RequestDiagnostics
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAgentDockerContainerListResponse parses an HTTP response from a AgentDockerContainerListWithResponse call
func ParseAgentDockerContainerListResponse(rsp *http.Response) (*AgentDockerContainerListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AgentDockerContainerListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []DockerContainer
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAgentDockerContainerInspectResponse parses an HTTP response from a AgentDockerContainerInspectWithResponse call
func ParseAgentDockerContainerInspectResponse(rsp *http.Response) (*AgentDockerContainerInspectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AgentDockerContainerInspectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []DockerContainerJSON
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAgentDockerImageListResponse parses an HTTP response from a AgentDockerImageListWithResponse call
func ParseAgentDockerImageListResponse(rsp *http.Response) (*AgentDockerImageListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AgentDockerImageListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []DockerImageSummary
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetAgentDockerInfoResponse parses an HTTP response from a GetAgentDockerInfoWithResponse call
func ParseGetAgentDockerInfoResponse(rsp *http.Response) (*GetAgentDockerInfoResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAgentDockerInfoResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest DockerInfo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAgentPullImageResponse parses an HTTP response from a AgentPullImageWithResponse call
func ParseAgentPullImageResponse(rsp *http.Response) (*AgentPullImageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AgentPullImageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
