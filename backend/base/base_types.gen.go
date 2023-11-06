// Package base provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package base

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for NamespaceKind.
const (
	NamespaceKindGroup            NamespaceKind = "group"
	NamespaceKindIntermediateCA   NamespaceKind = "int-ca"
	NamespaceKindProfile          NamespaceKind = "profile"
	NamespaceKindRootCA           NamespaceKind = "root-ca"
	NamespaceKindServicePrincipal NamespaceKind = "service-principal"
	NamespaceKindSystem           NamespaceKind = "sys"
	NamespaceKindUser             NamespaceKind = "user"
)

// Defines values for ResourceKind.
const (
	ProfileResourceKindGroup            ResourceKind = "group"
	ProfileResourceKindIntermediateCA   ResourceKind = "int-ca"
	ProfileResourceKindManagedApp       ResourceKind = "managed-app"
	ProfileResourceKindRootCA           ResourceKind = "root-ca"
	ProfileResourceKindServicePrincipal ResourceKind = "service-principal"
	ProfileResourceKindUser             ResourceKind = "user"
	ResourceKindAgentInstance           ResourceKind = "agent-instance"
	ResourceKindCert                    ResourceKind = "cert"
	ResourceKindCertPolicy              ResourceKind = "cert-policy"
	ResourceKindKeyPolicy               ResourceKind = "key-policy"
	ResourceKindNamespaceConfig         ResourceKind = "ns-config"
	ResourceKindSecretPolicy            ResourceKind = "secret-policy"
)

// AzureRoleAssignment defines model for AzureRoleAssignment.
type AzureRoleAssignment struct {
	ID               *string `json:"id,omitempty"`
	Name             *string `json:"name,omitempty"`
	PrincipalId      *string `json:"principalId,omitempty"`
	RoleDefinitionId *string `json:"roleDefinitionId,omitempty"`
}

// Base64RawURLEncodedBytes defines model for Base64RawURLEncodedBytes.
type Base64RawURLEncodedBytes = cloudkey.Base64RawURLEncodableBytes

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
	BuildID     string   `json:"buildId"`
	Environment []string `json:"environment,omitempty"`
	GoVersion   string   `json:"goVersion"`
}

// NamespaceIdentifierParameter defines model for NamespaceIdentifierParameter.
type NamespaceIdentifierParameter = Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = NamespaceKind

// ResourceIdentifierParameter defines model for ResourceIdentifierParameter.
type ResourceIdentifierParameter = Identifier

// AzureRoleAssignmentResponse defines model for AzureRoleAssignmentResponse.
type AzureRoleAssignmentResponse = AzureRoleAssignment

// ListAzureRoleAssignmentsResponse defines model for ListAzureRoleAssignmentsResponse.
type ListAzureRoleAssignmentsResponse = []AzureRoleAssignment
