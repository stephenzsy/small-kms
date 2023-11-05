// Package secret provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package secret

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for SecretGenerateMode.
const (
	SecretGenerateModeManual                SecretGenerateMode = "manual"
	SecretGenerateModeServerGeneratedRandom SecretGenerateMode = "random-server"
)

// Defines values for SecretRandomCharacterClass.
const (
	SecretRandomCharClassBase64RawURL SecretRandomCharacterClass = "base64-raw-url"
)

// SecretGenerateMode defines model for SecretGenerateMode.
type SecretGenerateMode string

// SecretPolicy defines model for SecretPolicy.
type SecretPolicy = secretPolicyComposed

// SecretPolicyFields defines model for SecretPolicyFields.
type SecretPolicyFields struct {
	ExpiryTime           *externalRef0.Period        `json:"expiryTime,omitempty"`
	Mode                 SecretGenerateMode          `json:"mode"`
	RandomCharacterClass *SecretRandomCharacterClass `json:"randomCharacterClass,omitempty"`

	// RandomLength Length of encoded random secret, in bytes
	RandomLength *int `json:"randomLength,omitempty"`
}

// SecretPolicyParameters defines model for SecretPolicyParameters.
type SecretPolicyParameters struct {
	DisplayName          string                      `json:"displayName,omitempty"`
	ExpiryTime           *externalRef0.Period        `json:"expiryTime,omitempty"`
	Mode                 SecretGenerateMode          `json:"mode"`
	RandomCharacterClass *SecretRandomCharacterClass `json:"randomCharacterClass,omitempty"`
	RandomLength         *int                        `json:"randomLength,omitempty"`
}

// SecretPolicyRef defines model for SecretPolicyRef.
type SecretPolicyRef = secretPolicyRefComposed

// SecretPolicyRefFields defines model for SecretPolicyRefFields.
type SecretPolicyRefFields struct {
	DisplayName string `json:"displayName"`
}

// SecretRandomCharacterClass defines model for SecretRandomCharacterClass.
type SecretRandomCharacterClass string

// SecretPolicyResponse defines model for SecretPolicyResponse.
type SecretPolicyResponse = SecretPolicy

// PutSecretPolicyJSONRequestBody defines body for PutSecretPolicy for application/json ContentType.
type PutSecretPolicyJSONRequestBody = SecretPolicyParameters

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List secret policies
	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/secret-policies)
	ListSecretPolicies(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter) error
	// Get key spec
	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/secret-policies/{resourceIdentifier})
	GetSecretPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter) error
	// Put key spec
	// (PUT /v1/{namespaceKind}/{namespaceIdentifier}/secret-policies/{resourceIdentifier})
	PutSecretPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListSecretPolicies converts echo context to params.
func (w *ServerInterfaceWrapper) ListSecretPolicies(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListSecretPolicies(ctx, namespaceKind, namespaceIdentifier)
	return err
}

// GetSecretPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) GetSecretPolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetSecretPolicy(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier)
	return err
}

// PutSecretPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) PutSecretPolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	// ------------- Path parameter "resourceIdentifier" -------------
	var resourceIdentifier externalRef0.ResourceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceIdentifier", runtime.ParamLocationPath, ctx.Param("resourceIdentifier"), &resourceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutSecretPolicy(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/secret-policies", wrapper.ListSecretPolicies)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/secret-policies/:resourceIdentifier", wrapper.GetSecretPolicy)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/secret-policies/:resourceIdentifier", wrapper.PutSecretPolicy)

}
