// Package profile provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package profile

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

// Profile defines model for Profile.
type Profile = ProfileRef

// ProfileParameters defines model for ProfileParameters.
type ProfileParameters struct {
	DisplayName *string `json:"displayName,omitempty"`
}

// ProfileRef defines model for ProfileRef.
type ProfileRef = profileRefComposed

// ProfileRefFields defines model for ProfileRefFields.
type ProfileRefFields struct {
	DisplayName string `json:"displayName"`
}

// ProfileResponse defines model for ProfileResponse.
type ProfileResponse = Profile

// PutRootCAJSONRequestBody defines body for PutRootCA for application/json ContentType.
type PutRootCAJSONRequestBody = ProfileParameters

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List root CA profiles
	// (GET /v1/root-ca)
	ListRootCAs(ctx echo.Context) error
	// Get profile
	// (GET /v1/root-ca/{namespaceIdentifier})
	GetRootCA(ctx echo.Context, namespaceIdentifier externalRef0.NamespaceIdentifierParameter) error
	// Put profile
	// (PUT /v1/root-ca/{namespaceIdentifier})
	PutRootCA(ctx echo.Context, namespaceIdentifier externalRef0.NamespaceIdentifierParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListRootCAs converts echo context to params.
func (w *ServerInterfaceWrapper) ListRootCAs(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListRootCAs(ctx)
	return err
}

// GetRootCA converts echo context to params.
func (w *ServerInterfaceWrapper) GetRootCA(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetRootCA(ctx, namespaceIdentifier)
	return err
}

// PutRootCA converts echo context to params.
func (w *ServerInterfaceWrapper) PutRootCA(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutRootCA(ctx, namespaceIdentifier)
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

	router.GET(baseURL+"/v1/root-ca", wrapper.ListRootCAs)
	router.GET(baseURL+"/v1/root-ca/:namespaceIdentifier", wrapper.GetRootCA)
	router.PUT(baseURL+"/v1/root-ca/:namespaceIdentifier", wrapper.PutRootCA)

}