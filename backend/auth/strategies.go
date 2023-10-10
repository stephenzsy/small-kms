package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func AuthorizeAdminOnly(c context.Context) error {
	if identity, ok := GetAuthIdentity(c); ok {
		if identity.HasAdminRole() {
			return nil
		}
	}
	return fmt.Errorf("%w: admin access required", common.ErrStatusForbidden)
}

func AuthorizeAgent(c context.Context) (uuid.UUID, error) {
	if identity, ok := GetAuthIdentity(c); ok {
		if identity.HasAppRole(roleKeyAgentActiveHost) {
			return identity.ClientPrincipalID(), nil
		}
	}
	return uuid.UUID{}, fmt.Errorf("%w: admin access required", common.ErrStatusForbidden)
}
