package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func AuthorizeAdminOnly(c context.Context) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasAdminRole()
	}
	return false
}

func AuthorizeSelfOrAdmin(c context.Context, namespaceID uuid.UUID) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasAdminRole() || (!utils.IsUUIDNil(namespaceID) && identity.ClientPrincipalID() == namespaceID)
	}
	return false
}

func HasRole(c context.Context, roleValue string) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasRole(roleValue)
	}
	return false
}

func ResolveSelfNamespace(c context.Context, nsUUID uuid.UUID, nsName string) uuid.UUID {
	if !utils.IsUUIDNil(nsUUID) || (!strings.EqualFold(nsName, "me") && !strings.EqualFold(nsName, "self")) {
		return nsUUID
	}
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.ClientPrincipalID()
	}
	return uuid.UUID{}
}
