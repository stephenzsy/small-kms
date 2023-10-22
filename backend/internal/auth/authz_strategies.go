package auth

import "context"

func AuthorizeAdminOnly(c context.Context) bool {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		return identity.HasAdminRole()
	}
	return false
}
