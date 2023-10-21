// Package base provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package base

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for NamespaceKind.
const (
	NamespaceKindProfile NamespaceKind = "profile"
)

// Defines values for ResourceKind.
const (
	ResourceKindManagedApp ResourceKind = "managed-app"
)

// Identifier defines model for Identifier.
type Identifier = identifierImpl

// NamespaceKind defines model for NamespaceKind.
type NamespaceKind string

// ResourceKind defines model for ResourceKind.
type ResourceKind string

// ResourceMetadata defines model for ResourceMetadata.
type ResourceMetadata struct {
	Deleted   *time.Time `json:"deleted,omitempty"`
	Updated   time.Time  `json:"updated"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// ResourceReference defines model for ResourceReference.
type ResourceReference struct {
	ID                  openapi_types.UUID `json:"id"`
	Metadata            ResourceMetadata   `json:"metadata"`
	NamespaceID         openapi_types.UUID `json:"namespaceId"`
	NamespaceIdentifier Identifier         `json:"namespaceIdentifier"`
	NamespaceKind       NamespaceKind      `json:"namespaceKind"`
	ResourceIdentifier  Identifier         `json:"resourceIdentifier"`
	ResourceKind        ResourceKind       `json:"resourceKind"`
}
