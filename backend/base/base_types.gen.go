// Package base provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package base

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for NamespaceKind.
const (
	NamespaceKindIntermediateCA   NamespaceKind = "int-ca"
	NamespaceKindProfile          NamespaceKind = "profile"
	NamespaceKindRootCA           NamespaceKind = "root-ca"
	NamespaceKindServicePrincipal NamespaceKind = "service-principal"
	NamespaceKindSystem           NamespaceKind = "sys"
)

// Defines values for ResourceKind.
const (
	ProfileResourceKindIntermediateCA   ResourceKind = "int-ca"
	ProfileResourceKindManagedApp       ResourceKind = "managed-app"
	ProfileResourceKindRootCA           ResourceKind = "root-ca"
	ProfileResourceKindServicePrincipal ResourceKind = "service-principal"
	ResourceKindCert                    ResourceKind = "cert"
	ResourceKindCertPolicy              ResourceKind = "cert-policy"
	ResourceKindKeyPolicy               ResourceKind = "key-policy"
	ResourceKindNamespaceConfig         ResourceKind = "ns-config"
)

// Base64RawURLEncodedBytes defines model for Base64RawURLEncodedBytes.
type Base64RawURLEncodedBytes = base64RawURLEncodedBytesImpl

// Identifier defines model for Identifier.
type Identifier = identifier

// NamespaceKind defines model for NamespaceKind.
type NamespaceKind string

// NumericDate defines model for NumericDate.
type NumericDate = jwt.NumericDate

// Period defines model for Period.
type Period = periodImpl

// RequestDiagnostics defines model for RequestDiagnostics.
type RequestDiagnostics struct {
	RequestHeaders []RequestHeaderEntry `json:"requestHeaders"`
	ServiceRuntime ServiceRuntimeInfo   `json:"serviceRuntime"`
}

// RequestHeaderEntry defines model for RequestHeaderEntry.
type RequestHeaderEntry struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// ResourceKind defines model for ResourceKind.
type ResourceKind string

// ResourceReference defines model for ResourceReference.
type ResourceReference struct {
	Deleted   *time.Time `json:"deleted,omitempty"`
	Id        Identifier `json:"id"`
	Updated   time.Time  `json:"updated"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// ResourceUniqueIdentifier defines model for ResourceUniqueIdentifier.
type ResourceUniqueIdentifier = DocFullIdentifier

// ServiceRuntimeInfo defines model for ServiceRuntimeInfo.
type ServiceRuntimeInfo struct {
	BuildID   string `json:"buildId"`
	GoVersion string `json:"goVersion"`
}

// NamespaceIdentifierParameter defines model for NamespaceIdentifierParameter.
type NamespaceIdentifierParameter = Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = NamespaceKind

// ResourceIdentifierParameter defines model for ResourceIdentifierParameter.
type ResourceIdentifierParameter = Identifier
