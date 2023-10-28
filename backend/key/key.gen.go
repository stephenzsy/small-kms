// Package key provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package key

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

// Defines values for JsonWebKeyCurveName.
const (
	JsonWebKeyCurveNameP256 JsonWebKeyCurveName = "P-256"
	JsonWebKeyCurveNameP384 JsonWebKeyCurveName = "P-384"
	JsonWebKeyCurveNameP521 JsonWebKeyCurveName = "P-521"
)

// Defines values for JsonWebKeyOperation.
const (
	JsonWebKeyOperationDecrypt   JsonWebKeyOperation = "decrypt"
	JsonWebKeyOperationEncrypt   JsonWebKeyOperation = "encrypt"
	JsonWebKeyOperationSign      JsonWebKeyOperation = "sign"
	JsonWebKeyOperationUnwrapKey JsonWebKeyOperation = "unwrapKey"
	JsonWebKeyOperationVerify    JsonWebKeyOperation = "verify"
	JsonWebKeyOperationWrapKey   JsonWebKeyOperation = "wrapKey"
)

// Defines values for JsonWebKeySignatureAlgorithm.
const (
	JsonWebKeySignatureAlgorithmES256 JsonWebKeySignatureAlgorithm = "ES256"
	JsonWebKeySignatureAlgorithmES384 JsonWebKeySignatureAlgorithm = "ES384"
	JsonWebKeySignatureAlgorithmES512 JsonWebKeySignatureAlgorithm = "ES512"
	JsonWebKeySignatureAlgorithmPS256 JsonWebKeySignatureAlgorithm = "PS256"
	JsonWebKeySignatureAlgorithmPS384 JsonWebKeySignatureAlgorithm = "PS384"
	JsonWebKeySignatureAlgorithmPS512 JsonWebKeySignatureAlgorithm = "PS512"
	JsonWebKeySignatureAlgorithmRS256 JsonWebKeySignatureAlgorithm = "RS256"
	JsonWebKeySignatureAlgorithmRS384 JsonWebKeySignatureAlgorithm = "RS384"
	JsonWebKeySignatureAlgorithmRS512 JsonWebKeySignatureAlgorithm = "RS512"
)

// Defines values for JsonWebKeyType.
const (
	JsonWebKeyTypeEC  JsonWebKeyType = "EC"
	JsonWebKeyTypeOct JsonWebKeyType = "oct"
	JsonWebKeyTypeRSA JsonWebKeyType = "RSA"
)

// JsonWebKey defines model for JsonWebKey.
type JsonWebKey struct {
	E             externalRef0.Base64RawURLEncodedBytes `json:"e,omitempty"`
	KeyOperations *[]JsonWebKeyOperation                `json:"key_ops,omitempty"`
	KeyID         *string                               `json:"kid,omitempty"`
	Kty           JsonWebKeyType                        `json:"kty"`
	N             externalRef0.Base64RawURLEncodedBytes `json:"n,omitempty"`
}

// JsonWebKeyCurveName defines model for JsonWebKeyCurveName.
type JsonWebKeyCurveName string

// JsonWebKeyOperation defines model for JsonWebKeyOperation.
type JsonWebKeyOperation string

// JsonWebKeySignatureAlgorithm defines model for JsonWebKeySignatureAlgorithm.
type JsonWebKeySignatureAlgorithm string

// JsonWebKeyType defines model for JsonWebKeyType.
type JsonWebKeyType string

// JsonWeyKeySize defines model for JsonWeyKeySize.
type JsonWeyKeySize = int32

// Key defines model for Key.
type Key = keyComposed

// KeyAttributes these attributes are not in JWK (RFC7517), more like JWT (RFC7519) fields
type KeyAttributes struct {
	Exp *externalRef0.NumericDate `json:"exp,omitempty"`
	Iat *externalRef0.NumericDate `json:"iat,omitempty"`
	Nbf *externalRef0.NumericDate `json:"nbf,omitempty"`
}

// KeyFields defines model for KeyFields.
type KeyFields struct {
	// Attributes these attributes are not in JWK (RFC7517), more like JWT (RFC7519) fields
	Attributes KeyAttributes `json:"attributes"`
}

// KeyPolicy defines model for KeyPolicy.
type KeyPolicy = keyPolicyComposed

// KeyPolicyFields defines model for KeyPolicyFields.
type KeyPolicyFields struct {
	ExpiryTime      *externalRef0.Period `json:"expiryTime,omitempty"`
	Exportable      bool                 `json:"exportable"`
	LifetimeActions []LifetimeAction     `json:"lifetimeActions,omitempty"`
}

// KeyPolicyParameters defines model for KeyPolicyParameters.
type KeyPolicyParameters struct {
	DisplayName *string              `json:"displayName,omitempty"`
	ExpiryTime  *externalRef0.Period `json:"expiryTime,omitempty"`
	Exportable  *bool                `json:"exportable,omitempty"`

	// KeySpec these attributes should mostly confirm to JWK (RFC7517)
	KeySpec         KeySpec          `json:"keySpec"`
	LifetimeActions []LifetimeAction `json:"lifetimeActions,omitempty"`
}

// KeyPolicyRef defines model for KeyPolicyRef.
type KeyPolicyRef = keyPolicyRefComposed

// KeyPolicyRefFields defines model for KeyPolicyRefFields.
type KeyPolicyRefFields struct {
	DisplayName string `json:"displayName"`
}

// KeySpec these attributes should mostly confirm to JWK (RFC7517)
type KeySpec struct {
	Crv           *JsonWebKeyCurveName                  `json:"crv,omitempty"`
	E             externalRef0.Base64RawURLEncodedBytes `json:"e,omitempty"`
	KeyOperations []JsonWebKeyOperation                 `json:"key_ops"`
	KeySize       *JsonWeyKeySize                       `json:"key_size,omitempty"`
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
	Alg           *JsonWebKeySignatureAlgorithm         `json:"alg,omitempty"`
	Crv           *JsonWebKeyCurveName                  `json:"crv,omitempty"`
	E             externalRef0.Base64RawURLEncodedBytes `json:"e,omitempty"`
	KeyOperations []JsonWebKeyOperation                 `json:"key_ops"`
	KeySize       *JsonWeyKeySize                       `json:"key_size,omitempty"`
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

// PutKeyPolicyJSONRequestBody defines body for PutKeyPolicy for application/json ContentType.
type PutKeyPolicyJSONRequestBody = KeyPolicyParameters

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List key policies
	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/key-policy)
	ListKeyPolicies(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter) error
	// Get key spec
	// (GET /v1/{namespaceKind}/{namespaceIdentifier}/key-policy/{resourceIdentifier})
	GetKeyPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter) error
	// Put key spec
	// (PUT /v1/{namespaceKind}/{namespaceIdentifier}/key-policy/{resourceIdentifier})
	PutKeyPolicy(ctx echo.Context, namespaceKind externalRef0.NamespaceKindParameter, namespaceIdentifier externalRef0.NamespaceIdentifierParameter, resourceIdentifier externalRef0.ResourceIdentifierParameter) error
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

	// ------------- Path parameter "namespaceIdentifier" -------------
	var namespaceIdentifier externalRef0.NamespaceIdentifierParameter

	err = runtime.BindStyledParameterWithLocation("simple", false, "namespaceIdentifier", runtime.ParamLocationPath, ctx.Param("namespaceIdentifier"), &namespaceIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter namespaceIdentifier: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeyPolicies(ctx, namespaceKind, namespaceIdentifier)
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
	err = w.Handler.GetKeyPolicy(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier)
	return err
}

// PutKeyPolicy converts echo context to params.
func (w *ServerInterfaceWrapper) PutKeyPolicy(ctx echo.Context) error {
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
	err = w.Handler.PutKeyPolicy(ctx, namespaceKind, namespaceIdentifier, resourceIdentifier)
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

	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/key-policy", wrapper.ListKeyPolicies)
	router.GET(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/key-policy/:resourceIdentifier", wrapper.GetKeyPolicy)
	router.PUT(baseURL+"/v1/:namespaceKind/:namespaceIdentifier/key-policy/:resourceIdentifier", wrapper.PutKeyPolicy)

}
