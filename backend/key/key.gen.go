// Package key provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package key

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// GenerateJsonWebKeyProperties defines model for GenerateJsonWebKeyProperties.
type GenerateJsonWebKeyProperties struct {
	Crv           JsonWebKeyCurveName   `json:"crv,omitempty"`
	KeyOperations []JsonWebKeyOperation `json:"key_ops,omitempty"`
	KeySize       *JsonWebKeySize       `json:"key_size,omitempty"`
	Kty           JsonWebKeyType        `json:"kty,omitempty"`
}

// JsonWebKey defines model for JsonWebKey.
type JsonWebKey = cloudkey.JsonWebKey

// JsonWebKeyCurveName defines model for JsonWebKeyCurveName.
type JsonWebKeyCurveName = cloudkey.JsonWebKeyCurveName

// JsonWebKeyOperation defines model for JsonWebKeyOperation.
type JsonWebKeyOperation = cloudkey.JsonWebKeyOperation

// JsonWebKeySize defines model for JsonWebKeySize.
type JsonWebKeySize = int32

// JsonWebKeyType defines model for JsonWebKeyType.
type JsonWebKeyType = cloudkey.JsonWebKeyType

// JsonWebSignatureAlgorithm defines model for JsonWebSignatureAlgorithm.
type JsonWebSignatureAlgorithm = cloudkey.JsonWebSignatureAlgorithm

// JsonWebSignatureKey defines model for JsonWebSignatureKey.
type JsonWebSignatureKey = cloudkey.JsonWebKey

// Key defines model for Key.
type Key = keyComposed

// KeyFields defines model for KeyFields.
type KeyFields struct {
	Iat     externalRef0.NumericDate     `json:"iat"`
	KeySize *JsonWebKeySize              `json:"key_size,omitempty"`
	Nbf     *externalRef0.NumericDate    `json:"nbf,omitempty"`
	Policy  externalRef0.ResourceLocator `json:"policy"`
}

// KeyPolicy defines model for KeyPolicy.
type KeyPolicy = keyPolicyComposed

// KeyPolicyFields defines model for KeyPolicyFields.
type KeyPolicyFields struct {
	ExpiryTime    *externalRef0.Period         `json:"expiryTime,omitempty"`
	Exportable    bool                         `json:"exportable"`
	KeyProperties GenerateJsonWebKeyProperties `json:"keyProperties"`
}

// KeyPolicyParameters defines model for KeyPolicyParameters.
type KeyPolicyParameters struct {
	DisplayName   string                        `json:"displayName,omitempty"`
	ExpiryTime    *externalRef0.Period          `json:"expiryTime,omitempty"`
	Exportable    *bool                         `json:"exportable,omitempty"`
	KeyProperties *GenerateJsonWebKeyProperties `json:"keyProperties,omitempty"`
}

// KeyPolicyRef defines model for KeyPolicyRef.
type KeyPolicyRef = keyPolicyRefComposed

// KeyPolicyRefFields defines model for KeyPolicyRefFields.
type KeyPolicyRefFields struct {
	DisplayName string `json:"displayName"`
}

// KeyRef defines model for KeyRef.
type KeyRef = keyRefComposed

// KeyRefFields defines model for KeyRefFields.
type KeyRefFields struct {
	Exp *externalRef0.NumericDate `json:"exp,omitempty"`
	Iat externalRef0.NumericDate  `json:"iat"`
}

// KeySpec these attributes should mostly confirm to JWK (RFC7517)
type KeySpec struct {
	Crv           JsonWebKeyCurveName                   `json:"crv,omitempty"`
	E             externalRef0.Base64RawURLEncodedBytes `json:"e,omitempty"`
	KeyOperations []JsonWebKeyOperation                 `json:"key_ops"`
	KeySize       *JsonWebKeySize                       `json:"key_size,omitempty"`
	KeyID         *string                               `json:"kid,omitempty"`
	Kty           JsonWebKeyType                        `json:"kty"`
	N             externalRef0.Base64RawURLEncodedBytes `json:"n,omitempty"`
	X             externalRef0.Base64RawURLEncodedBytes `json:"x,omitempty"`
	Y             externalRef0.Base64RawURLEncodedBytes `json:"y,omitempty"`
}

// LifetimeAction defines model for LifetimeAction.
type LifetimeAction struct {
	Trigger LifetimeTrigger `json:"trigger"`
}

// LifetimeTrigger defines model for LifetimeTrigger.
type LifetimeTrigger struct {
	PercentageAfterCreate *int32               `json:"percentageAfterCreate,omitempty"`
	TimeAfterCreate       *externalRef0.Period `json:"timeAfterCreate,omitempty"`
	TimeBeforeExpiry      *externalRef0.Period `json:"timeBeforeExpiry,omitempty"`
}

// SigningKeySpec defines model for SigningKeySpec.
type SigningKeySpec struct {
	Alg           *JsonWebSignatureAlgorithm            `json:"alg,omitempty"`
	Crv           JsonWebKeyCurveName                   `json:"crv,omitempty"`
	E             externalRef0.Base64RawURLEncodedBytes `json:"e,omitempty"`
	KeyOperations []JsonWebKeyOperation                 `json:"key_ops"`
	KeySize       *JsonWebKeySize                       `json:"key_size,omitempty"`
	KeyID         *string                               `json:"kid,omitempty"`
	Kty           JsonWebKeyType                        `json:"kty"`
	N             externalRef0.Base64RawURLEncodedBytes `json:"n,omitempty"`
	X             externalRef0.Base64RawURLEncodedBytes `json:"x,omitempty"`

	// X5c Base64 encoded certificate chain
	CertificateChain []externalRef0.Base64RawURLEncodedBytes `json:"x5c,omitempty"`
	X5t              externalRef0.Base64RawURLEncodedBytes   `json:"x5t,omitempty"`
	X5tS256          externalRef0.Base64RawURLEncodedBytes   `json:"x5t#S256,omitempty"`
	Y                externalRef0.Base64RawURLEncodedBytes   `json:"y,omitempty"`
}

// KeyPolicyResponse defines model for KeyPolicyResponse.
type KeyPolicyResponse = KeyPolicy

// KeyResponse defines model for KeyResponse.
type KeyResponse = Key

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List key policies
	// (GET /v1/{namespaceKind}/{namespaceId}/key-policies)
	ListKeyPolicies(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter) error
	// Get key spec
	// (GET /v1/{namespaceKind}/{namespaceId}/key-policies/{resourceId})
	GetKeyPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceId externalRef0.NamespaceIdParameter, resourceId externalRef0.ResourceIdParameter) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListKeyPolicies converts echo context to params.
func (w *ServerInterfaceWrapper) ListKeyPolicies(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeyPolicies(ctx, namespaceKind, namespaceId)
	return err
}

// GetKeyPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) GetKeyPolicy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "namespaceKind" -------------
	var namespaceKind externalRef0.NamespaceKindParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceKind", runtime.ParamLocationPath, ctx.Param("namespaceKind"), &namespaceKind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceKind: %s", err))
	}

	// ------------- Path parameter "namespaceId" -------------
	var namespaceId externalRef0.NamespaceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceId", runtime.ParamLocationPath, ctx.Param("namespaceId"), &namespaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceId: %s", err))
	}

	// ------------- Path parameter "resourceId" -------------
	var resourceId externalRef0.ResourceIdParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, ctx.Param("resourceId"), &resourceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resourceId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetKeyPolicy(ctx, namespaceKind, namespaceId, resourceId)
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

	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/key-policies", wrapper.ListKeyPolicies)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceId/key-policies/:resourceId", wrapper.GetKeyPolicy)

}
