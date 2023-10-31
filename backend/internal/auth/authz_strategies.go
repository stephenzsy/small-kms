package auth

import (
	"context"

	"github.com/google/uuid"
)

func AuthorizeAdminOnly(c context.Context) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasAdminRole()
	}
	return false
}

func AuthorizeApplicationOrAdmin(c context.Context, namespaceID uuid.UUID) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasAdminRole() || identity.ClientPrincipalID() == namespaceID
	}
	return false
}

func AuthorizeApplicationMe(c context.Context, namespaceID uuid.UUID, me bool) (uuid.UUID, bool) {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.ClientPrincipalID(), me || identity.ClientPrincipalID() == namespaceID
	}
	return uuid.UUID{}, false
}

func HasRole(c context.Context, roleValue string) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasRole(roleValue)
	}
	return false
}
