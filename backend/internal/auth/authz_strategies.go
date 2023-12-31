package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// Deprecated use authz.AuthorizeAdminOnly instead.
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

func ResolveSelfNamespace(c context.Context, nsID string) uuid.UUID {
	isUUID := false
	var nsUUID uuid.UUID
	var err error
	if nsUUID, err = uuid.Parse(nsID); err == nil {
		isUUID = true
	}
	if (isUUID && nsUUID != uuid.UUID{}) || (!strings.EqualFold(nsID, "me") && !strings.EqualFold(nsID, "self")) {
		return nsUUID
	}
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.ClientPrincipalID()
	}
	return uuid.UUID{}
}
