// Package agentclient provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
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
	externalRef0 "github.com/stephenzsy/small-kms/backend/shared"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for IncludeCertificate.
const (
	IncludeJWK IncludeCertificate = "jwk"
	IncludePEM IncludeCertificate = "pem"
)

// IncludeCertificate defines model for IncludeCertificate.
type IncludeCertificate string

// AgentConfigNameParameter defines model for AgentConfigNameParameter.
type AgentConfigNameParameter = externalRef0.AgentConfigName

// CertificateIdPathParameter defines model for CertificateIdPathParameter.
type CertificateIdPathParameter = externalRef0.Identifier

// CertificateTemplateIdentifierParameter defines model for CertificateTemplateIdentifierParameter.
type CertificateTemplateIdentifierParameter = externalRef0.Identifier

// IncludeCertificateParameter defines model for IncludeCertificateParameter.
type IncludeCertificateParameter = IncludeCertificate

// NamespaceIdParameter defines model for NamespaceIdParameter.
type NamespaceIdParameter = externalRef0.Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = externalRef0.NamespaceKind

// AgentConfigurationResponse defines model for AgentConfigurationResponse.
type AgentConfigurationResponse = externalRef0.AgentConfiguration

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = externalRef0.CertificateInfo

// GetAgentConfigurationParams defines parameters for GetAgentConfiguration.
type GetAgentConfigurationParams struct {
	RefreshToken               *string `form:"refreshToken,omitempty" json:"refreshToken,omitempty"`
	XSmallkmsIfVersionNotMatch *string `json:"X-Smallkms-If-Version-Not-Match,omitempty"`
}

// GetCertificateParams defines parameters for GetCertificate.
type GetCertificateParams struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
	TemplateId         *externalRef0.Identifier     `form:"templateId,omitempty" json:"templateId,omitempty"`
}

// AgentCallbackJSONRequestBody defines body for AgentCallback for application/json ContentType.
type AgentCallbackJSONRequestBody = externalRef0.AgentCallbackRequest

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
	// GetDiagnostics request
	GetDiagnostics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AgentCallbackWithBody request with any body
	AgentCallbackWithBody(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AgentCallback(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, body AgentCallbackJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAgentConfiguration request
	GetAgentConfiguration(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params *GetAgentConfigurationParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCertificate request
	GetCertificate(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params *GetCertificateParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetDiagnostics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDiagnosticsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentCallbackWithBody(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentCallbackRequestWithBody(c.Server, namespaceKind, namespaceId, configName, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AgentCallback(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, body AgentCallbackJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAgentCallbackRequest(c.Server, namespaceKind, namespaceId, configName, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAgentConfiguration(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params *GetAgentConfigurationParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAgentConfigurationRequest(c.Server, namespaceKind, namespaceId, configName, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCertificate(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params *GetCertificateParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCertificateRequest(c.Server, namespaceKind, namespaceId, certificateId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetDiagnosticsRequest generates requests for GetDiagnostics
func NewGetDiagnosticsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v3/diagnostics")
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

// NewAgentCallbackRequest calls the generic AgentCallback builder with application/json body
func NewAgentCallbackRequest(server string, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, body AgentCallbackJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAgentCallbackRequestWithBody(server, namespaceKind, namespaceId, configName, "application/json", bodyReader)
}

// NewAgentCallbackRequestWithBody generates requests for AgentCallback with any type of body
func NewAgentCallbackRequestWithBody(server string, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, contentType string, body io.Reader) (*http.Request, error) {
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

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "configName", runtime.ParamLocationPath, configName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v3/%s/%s/agent-callback/%s", pathParam0, pathParam1, pathParam2)
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

	return req, nil
}

// NewGetAgentConfigurationRequest generates requests for GetAgentConfiguration
func NewGetAgentConfigurationRequest(server string, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params *GetAgentConfigurationParams) (*http.Request, error) {
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

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "configName", runtime.ParamLocationPath, configName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v3/%s/%s/agent-config/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.RefreshToken != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "refreshToken", runtime.ParamLocationQuery, *params.RefreshToken); err != nil {
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

	if params != nil {

		if params.XSmallkmsIfVersionNotMatch != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Smallkms-If-Version-Not-Match", runtime.ParamLocationHeader, *params.XSmallkmsIfVersionNotMatch)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Smallkms-If-Version-Not-Match", headerParam0)
		}

	}

	return req, nil
}

// NewGetCertificateRequest generates requests for GetCertificate
func NewGetCertificateRequest(server string, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params *GetCertificateParams) (*http.Request, error) {
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

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "certificateId", runtime.ParamLocationPath, certificateId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v3/%s/%s/certificate/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.IncludeCertificate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "includeCertificate", runtime.ParamLocationQuery, *params.IncludeCertificate); err != nil {
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

		if params.TemplateId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "templateId", runtime.ParamLocationQuery, *params.TemplateId); err != nil {
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
	// GetDiagnosticsWithResponse request
	GetDiagnosticsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetDiagnosticsResponse, error)

	// AgentCallbackWithBodyWithResponse request with any body
	AgentCallbackWithBodyWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AgentCallbackResponse, error)

	AgentCallbackWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, body AgentCallbackJSONRequestBody, reqEditors ...RequestEditorFn) (*AgentCallbackResponse, error)

	// GetAgentConfigurationWithResponse request
	GetAgentConfigurationWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params *GetAgentConfigurationParams, reqEditors ...RequestEditorFn) (*GetAgentConfigurationResponse, error)

	// GetCertificateWithResponse request
	GetCertificateWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params *GetCertificateParams, reqEditors ...RequestEditorFn) (*GetCertificateResponse, error)
}

type GetDiagnosticsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.RequestDiagnostics
}

// Status returns HTTPResponse.Status
func (r GetDiagnosticsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDiagnosticsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AgentCallbackResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r AgentCallbackResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AgentCallbackResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAgentConfigurationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AgentConfigurationResponse
}

// Status returns HTTPResponse.Status
func (r GetAgentConfigurationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAgentConfigurationResponse) StatusCode() int {
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

// GetDiagnosticsWithResponse request returning *GetDiagnosticsResponse
func (c *ClientWithResponses) GetDiagnosticsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetDiagnosticsResponse, error) {
	rsp, err := c.GetDiagnostics(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDiagnosticsResponse(rsp)
}

// AgentCallbackWithBodyWithResponse request with arbitrary body returning *AgentCallbackResponse
func (c *ClientWithResponses) AgentCallbackWithBodyWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AgentCallbackResponse, error) {
	rsp, err := c.AgentCallbackWithBody(ctx, namespaceKind, namespaceId, configName, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentCallbackResponse(rsp)
}

func (c *ClientWithResponses) AgentCallbackWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, body AgentCallbackJSONRequestBody, reqEditors ...RequestEditorFn) (*AgentCallbackResponse, error) {
	rsp, err := c.AgentCallback(ctx, namespaceKind, namespaceId, configName, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAgentCallbackResponse(rsp)
}

// GetAgentConfigurationWithResponse request returning *GetAgentConfigurationResponse
func (c *ClientWithResponses) GetAgentConfigurationWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, configName AgentConfigNameParameter, params *GetAgentConfigurationParams, reqEditors ...RequestEditorFn) (*GetAgentConfigurationResponse, error) {
	rsp, err := c.GetAgentConfiguration(ctx, namespaceKind, namespaceId, configName, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAgentConfigurationResponse(rsp)
}

// GetCertificateWithResponse request returning *GetCertificateResponse
func (c *ClientWithResponses) GetCertificateWithResponse(ctx context.Context, namespaceKind NamespaceKindParameter, namespaceId NamespaceIdParameter, certificateId CertificateIdPathParameter, params *GetCertificateParams, reqEditors ...RequestEditorFn) (*GetCertificateResponse, error) {
	rsp, err := c.GetCertificate(ctx, namespaceKind, namespaceId, certificateId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCertificateResponse(rsp)
}

// ParseGetDiagnosticsResponse parses an HTTP response from a GetDiagnosticsWithResponse call
func ParseGetDiagnosticsResponse(rsp *http.Response) (*GetDiagnosticsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDiagnosticsResponse{
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

// ParseAgentCallbackResponse parses an HTTP response from a AgentCallbackWithResponse call
func ParseAgentCallbackResponse(rsp *http.Response) (*AgentCallbackResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AgentCallbackResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetAgentConfigurationResponse parses an HTTP response from a GetAgentConfigurationWithResponse call
func ParseGetAgentConfigurationResponse(rsp *http.Response) (*GetAgentConfigurationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAgentConfigurationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AgentConfigurationResponse
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
