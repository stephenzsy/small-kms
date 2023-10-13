package ns

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func ResolveAuthedNamespaseID(c context.Context, namespaceKind shared.NamespaceKind, inputID shared.Identifier) (shared.Identifier, error) {
	if a, ok := auth.GetAuthIdentity(c); ok {
		if inputID.IsNilOrEmpty() || (!inputID.IsUUID() && inputID.String() == "me") {
			inputID = shared.UUIDIdentifier(a.ClientPrincipalID())
		} else if inputID.IsUUID() && inputID.UUID() == a.ClientPrincipalID() {
			// ok
		} else {
			return inputID, fmt.Errorf("%w: authroization namespace mismatch: %s", common.ErrStatusForbidden, a.ClientPrincipalID())
		}
		switch namespaceKind {
		case shared.NamespaceKindUser,
			shared.NamespaceKindServicePrincipal:
			// ok
		default:
			return inputID, fmt.Errorf("%w: invalid namespaceKind: %s", common.ErrStatusForbidden, namespaceKind)
		}
	}

	return inputID, fmt.Errorf("%w: no authorization context", common.ErrStatusUnauthorized)
}
