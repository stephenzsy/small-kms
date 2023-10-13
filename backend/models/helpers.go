package models

import (
	"github.com/stephenzsy/small-kms/backend/shared"
)

// Deprecated: use shared.NamespaceKind instead
type NamespaceKind = shared.NamespaceKind

// Deprecated: use shared.ResourceKind instead
type ResourceKind = shared.ResourceKind

// Deprecated: use shared.NamespaceIdentifier instead
type NamespaceID = shared.NamespaceIdentifier

// Deprecated: use shared.ResourceIdentifer instead
type ResourceID = shared.ResourceIdentifier

// Deprecated: use shared.ResourceLocator instead
type ResourceLocator = shared.ResourceLocator

// Deprecated: use shared.NewResourceLocator instead
func NewResourceLocator(namespaceID NamespaceID, resourceID ResourceID) shared.ResourceLocator {
	return shared.NewResourceLocator(namespaceID, resourceID)
}

// Deprecated: use shared.NewNamespaceIdentifier instead
func NewNamespaceID(kind NamespaceKind, identifier shared.Identifier) NamespaceID {
	return shared.NewNamespaceIdentifier(kind, identifier)
}
