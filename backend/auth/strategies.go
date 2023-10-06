package auth

import (
	"context"
	"fmt"

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
