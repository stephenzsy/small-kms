package models

import (
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func (alg JwkAlg) ToAzKeysSignatureAlgorithm() azkeys.SignatureAlgorithm {
	switch alg {
	case AlgRS256:
		return azkeys.SignatureAlgorithmRS256
	case AlgRS384:
		return azkeys.SignatureAlgorithmRS384
	case AlgRS512:
		return azkeys.SignatureAlgorithmRS512
	case AlgES256:
		return azkeys.SignatureAlgorithmES256
	case AlgES384:
		return azkeys.SignatureAlgorithmES384
	}
	return azkeys.SignatureAlgorithm("")
}

func (alg JwkAlg) ToX509SignatureAlgorithm() x509.SignatureAlgorithm {
	switch alg {
	case AlgRS256:
		return x509.SHA256WithRSA
	case AlgRS384:
		return x509.SHA384WithRSA
	case AlgRS512:
		return x509.SHA512WithRSA
	case AlgES256:
		return x509.ECDSAWithSHA256
	case AlgES384:
		return x509.ECDSAWithSHA384
	}
	return x509.UnknownSignatureAlgorithm
}

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

// Deprecated: use NamespaceKind instead
type ProfileType = NamespaceKind

// Deprecated: use shared.NewResourceLocator instead
func NewResourceLocator(namespaceID NamespaceID, resourceID ResourceID) shared.ResourceLocator {
	return shared.NewResourceLocator(namespaceID, resourceID)
}

// Deprecated: use shared.NewNamespaceIdentifier instead
func NewNamespaceID(kind NamespaceKind, identifier shared.Identifier) NamespaceID {
	return shared.NewNamespaceIdentifier(kind, identifier)
}

/*
	func NewNamespaceStringID(kind NamespaceKind, id string) NamespaceID {
		return common.NewIdentifierWithKind(kind, common.StringIdentifier(id))
	}
*/

// Deprecated: use shared.NewResourceIdentifier instead
func NewResourceID(kind ResourceKind, identifier shared.Identifier) ResourceID {
	return shared.NewResourceIdentifier(kind, identifier)
}
