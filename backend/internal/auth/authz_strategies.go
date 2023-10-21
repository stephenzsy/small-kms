package auth

import (
	"net/http"

	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func AuthorizeAdminOnly(c ctx.RequestContext) error {
	if identity, ok := c.Value(authIdentityContextKey).(AuthIdentity); ok {
		if identity.HasAdminRole() {
			return nil
		}
	}
	return c.JSON(http.StatusForbidden, map[string]string{"message": "admin access required"})
}
