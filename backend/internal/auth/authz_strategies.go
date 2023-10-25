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
