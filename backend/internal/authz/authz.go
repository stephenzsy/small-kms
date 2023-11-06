package authz

import ctx "github.com/stephenzsy/small-kms/backend/internal/context"

type AuthzResult uint

type RequestContext = ctx.RequestContext

const (
	AuthzResultNone AuthzResult = iota
	AuthzResultAllow
	AuthzResultDeny
)

type AuthZFunc = func(c RequestContext) (RequestContext, AuthzResult)

func Authorize(c RequestContext, authzFns ...AuthZFunc) (RequestContext, bool) {
	var r AuthzResult
	for _, authzFn := range authzFns {
		if authzFn == nil {
			continue
		}
		if c, r = authzFn(c); r == AuthzResultDeny {
			return c, false
		} else if r == AuthzResultAllow {
			return c, true
		} // else to next handler
	}
	return c, false
}
