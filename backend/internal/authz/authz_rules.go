package authz

import "github.com/stephenzsy/small-kms/backend/internal/auth"

func AllowAdmin(c RequestContext) (RequestContext, AuthzResult) {
	identity := auth.GetAuthIdentity(c)
	if identity.HasAdminRole() {
		return c, AuthzResultAllow
	}
	return c, AuthzResultNone
}

var _ AuthZFunc = AllowAdmin

func AllowSelf(namespaceId string) AuthZFunc {
	return func(c RequestContext) (RequestContext, AuthzResult) {
		identity := auth.GetAuthIdentity(c)
		if identity.ClientPrincipalID().String() == namespaceId {
			return c, AuthzResultAllow
		}
		return c, AuthzResultNone
	}
}

var _ AuthZFunc = AllowAdmin

// convinient function to authorize admin only, context should not be modified
func AuthorizeAdminOnly(c RequestContext) bool {
	_, ok := Authorize(c, AllowAdmin)
	return ok
}
