// Package agentclient provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package agentclient

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
	externalRef1 "github.com/stephenzsy/small-kms/backend/cert"
	externalRef3 "github.com/stephenzsy/small-kms/backend/managedapp"
	externalRef4 "github.com/stephenzsy/small-kms/backend/secret"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AgentConfigServer defines model for AgentConfigServer.
type AgentConfigServer = externalRef3.AgentConfigServer

// AgentInstanceFields defines model for AgentInstanceFields.
type AgentInstanceFields = externalRef3.AgentInstanceFields

// Certificate defines model for Certificate.
type Certificate = externalRef1.Certificate

// EnrollCertificateRequest defines model for EnrollCertificateRequest.
type EnrollCertificateRequest = externalRef1.EnrollCertificateRequest

// Secret defines model for Secret.
type Secret = externalRef4.Secret

// AgentConfigRadiusResponse defines model for AgentConfigRadiusResponse.
type AgentConfigRadiusResponse = externalRef3.AgentConfigRadius

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = externalRef1.Certificate

// EnrollCertificateParams defines parameters for EnrollCertificate.
type EnrollCertificateParams struct {
	DryRun *bool `form:"dryRun,omitempty" json:"dryRun,omitempty"`
}

// GetSecretParams defines parameters for GetSecret.
type GetSecretParams struct {
	WithValue *bool `form:"withValue,omitempty" json:"withValue,omitempty"`
}

// PutAgentInstanceJSONRequestBody defines body for PutAgentInstance for application/json ContentType.
type PutAgentInstanceJSONRequestBody = AgentInstanceFields

// EnrollCertificateJSONRequestBody defines body for EnrollCertificate for application/json ContentType.
type EnrollCertificateJSONRequestBody = EnrollCertificateRequest

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
	// GetAgentConfigRadius request
	GetAgentConfigRadius(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAgentConfigServer request
	GetAgentConfigServer(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PutAgentInstanceWithBody request with any body
	PutAgentInstanceWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PutAgentInstance(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, body PutAgentInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// EnrollCertificateWithBody request with any body
	EnrollCertificateWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	EnrollCertificate(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, body EnrollCertificateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCertificate request
	GetCertificate(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSecret request
	GetSecret(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *GetSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAgentConfigRadius(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAgentConfigRadiusRequest(c.Server, namespaceKind, namespaceId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAgentConfigServer(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAgentConfigServerRequest(c.Server, namespaceKind, namespaceId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutAgentInstanceWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutAgentInstanceRequestWithBody(c.Server, namespaceKind, namespaceId, resourceId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutAgentInstance(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, body PutAgentInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutAgentInstanceRequest(c.Server, namespaceKind, namespaceId, resourceId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) EnrollCertificateWithBody(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewEnrollCertificateRequestWithBody(c.Server, namespaceKind, namespaceId, resourceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) EnrollCertificate(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, body EnrollCertificateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewEnrollCertificateRequest(c.Server, namespaceKind, namespaceId, resourceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCertificate(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCertificateRequest(c.Server, namespaceKind, namespaceId, resourceId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetSecret(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *GetSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSecretRequest(c.Server, namespaceKind, namespaceId, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAgentConfigRadiusRequest generates requests for GetAgentConfigRadius
func NewGetAgentConfigRadiusRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent-config/radius", pathParam0, pathParam1)
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

	return req, nil
}

// NewGetAgentConfigServerRequest generates requests for GetAgentConfigServer
func NewGetAgentConfigServerRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent-config/server", pathParam0, pathParam1)
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

	return req, nil
}

// NewPutAgentInstanceRequest calls the generic PutAgentInstance builder with application/json body
func NewPutAgentInstanceRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, body PutAgentInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPutAgentInstanceRequestWithBody(server, namespaceKind, namespaceId, resourceId, "application/json", bodyReader)
}

// NewPutAgentInstanceRequestWithBody generates requests for PutAgentInstance with any type of body
func NewPutAgentInstanceRequestWithBody(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/agent/instance/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewEnrollCertificateRequest calls the generic EnrollCertificate builder with application/json body
func NewEnrollCertificateRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, body EnrollCertificateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewEnrollCertificateRequestWithBody(server, namespaceKind, namespaceId, resourceId, params, "application/json", bodyReader)
}

// NewEnrollCertificateRequestWithBody generates requests for EnrollCertificate with any type of body
func NewEnrollCertificateRequestWithBody(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/cert-policy/%s/enroll-cert", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.DryRun != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "dryRun", runtime.ParamLocationQuery, *params.DryRun); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetCertificateRequest generates requests for GetCertificate
func NewGetCertificateRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/cert/%s", pathParam0, pathParam1, pathParam2)
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

	return req, nil
}

// NewGetSecretRequest generates requests for GetSecret
func NewGetSecretRequest(server string, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *GetSecretParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, namespaceKind)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, namespaceId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/%s/%s/secrets/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.WithValue != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "withValue", runtime.ParamLocationQuery, *params.WithValue); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
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
	// GetAgentConfigRadiusWithResponse request
	GetAgentConfigRadiusWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*GetAgentConfigRadiusResponse, error)

	// GetAgentConfigServerWithResponse request
	GetAgentConfigServerWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*GetAgentConfigServerResponse, error)

	// PutAgentInstanceWithBodyWithResponse request with any body
	PutAgentInstanceWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutAgentInstanceResponse, error)

	PutAgentInstanceWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, body PutAgentInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*PutAgentInstanceResponse, error)

	// EnrollCertificateWithBodyWithResponse request with any body
	EnrollCertificateWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*EnrollCertificateResponse, error)

	EnrollCertificateWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, body EnrollCertificateJSONRequestBody, reqEditors ...RequestEditorFn) (*EnrollCertificateResponse, error)

	// GetCertificateWithResponse request
	GetCertificateWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, reqEditors ...RequestEditorFn) (*GetCertificateResponse, error)

	// GetSecretWithResponse request
	GetSecretWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *GetSecretParams, reqEditors ...RequestEditorFn) (*GetSecretResponse, error)
}

type GetAgentConfigRadiusResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AgentConfigRadiusResponse
}

// Status returns HTTPResponse.Status
func (r GetAgentConfigRadiusResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAgentConfigRadiusResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAgentConfigServerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AgentConfigServer
}

// Status returns HTTPResponse.Status
func (r GetAgentConfigServerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAgentConfigServerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PutAgentInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r PutAgentInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PutAgentInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type EnrollCertificateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Certificate
}

// Status returns HTTPResponse.Status
func (r EnrollCertificateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r EnrollCertificateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCertificateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CertificateResponse
}

// Status returns HTTPResponse.Status
func (r GetCertificateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCertificateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Secret
}

// Status returns HTTPResponse.Status
func (r GetSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAgentConfigRadiusWithResponse request returning *GetAgentConfigRadiusResponse
func (c *ClientWithResponses) GetAgentConfigRadiusWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*GetAgentConfigRadiusResponse, error) {
	rsp, err := c.GetAgentConfigRadius(ctx, namespaceKind, namespaceId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAgentConfigRadiusResponse(rsp)
}

// GetAgentConfigServerWithResponse request returning *GetAgentConfigServerResponse
func (c *ClientWithResponses) GetAgentConfigServerWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, reqEditors ...RequestEditorFn) (*GetAgentConfigServerResponse, error) {
	rsp, err := c.GetAgentConfigServer(ctx, namespaceKind, namespaceId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAgentConfigServerResponse(rsp)
}

// PutAgentInstanceWithBodyWithResponse request with arbitrary body returning *PutAgentInstanceResponse
func (c *ClientWithResponses) PutAgentInstanceWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutAgentInstanceResponse, error) {
	rsp, err := c.PutAgentInstanceWithBody(ctx, namespaceKind, namespaceId, resourceId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutAgentInstanceResponse(rsp)
}

func (c *ClientWithResponses) PutAgentInstanceWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, body PutAgentInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*PutAgentInstanceResponse, error) {
	rsp, err := c.PutAgentInstance(ctx, namespaceKind, namespaceId, resourceId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutAgentInstanceResponse(rsp)
}

// EnrollCertificateWithBodyWithResponse request with arbitrary body returning *EnrollCertificateResponse
func (c *ClientWithResponses) EnrollCertificateWithBodyWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*EnrollCertificateResponse, error) {
	rsp, err := c.EnrollCertificateWithBody(ctx, namespaceKind, namespaceId, resourceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseEnrollCertificateResponse(rsp)
}

func (c *ClientWithResponses) EnrollCertificateWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *EnrollCertificateParams, body EnrollCertificateJSONRequestBody, reqEditors ...RequestEditorFn) (*EnrollCertificateResponse, error) {
	rsp, err := c.EnrollCertificate(ctx, namespaceKind, namespaceId, resourceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseEnrollCertificateResponse(rsp)
}

// GetCertificateWithResponse request returning *GetCertificateResponse
func (c *ClientWithResponses) GetCertificateWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, reqEditors ...RequestEditorFn) (*GetCertificateResponse, error) {
	rsp, err := c.GetCertificate(ctx, namespaceKind, namespaceId, resourceId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCertificateResponse(rsp)
}

// GetSecretWithResponse request returning *GetSecretResponse
func (c *ClientWithResponses) GetSecretWithResponse(ctx context.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter, params *GetSecretParams, reqEditors ...RequestEditorFn) (*GetSecretResponse, error) {
	rsp, err := c.GetSecret(ctx, namespaceKind, namespaceId, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSecretResponse(rsp)
}

// ParseGetAgentConfigRadiusResponse parses an HTTP response from a GetAgentConfigRadiusWithResponse call
func ParseGetAgentConfigRadiusResponse(rsp *http.Response) (*GetAgentConfigRadiusResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAgentConfigRadiusResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AgentConfigRadiusResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetAgentConfigServerResponse parses an HTTP response from a GetAgentConfigServerWithResponse call
func ParseGetAgentConfigServerResponse(rsp *http.Response) (*GetAgentConfigServerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAgentConfigServerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AgentConfigServer
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePutAgentInstanceResponse parses an HTTP response from a PutAgentInstanceWithResponse call
func ParsePutAgentInstanceResponse(rsp *http.Response) (*PutAgentInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PutAgentInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseEnrollCertificateResponse parses an HTTP response from a EnrollCertificateWithResponse call
func ParseEnrollCertificateResponse(rsp *http.Response) (*EnrollCertificateResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &EnrollCertificateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Certificate
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetCertificateResponse parses an HTTP response from a GetCertificateWithResponse call
func ParseGetCertificateResponse(rsp *http.Response) (*GetCertificateResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCertificateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CertificateResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetSecretResponse parses an HTTP response from a GetSecretWithResponse call
func ParseGetSecretResponse(rsp *http.Response) (*GetSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Secret
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
