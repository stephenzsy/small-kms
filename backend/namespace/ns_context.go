package ns

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type NSContext interface {
	Kind() base.NamespaceKind
	ID() base.ID
}

type nsContext struct {
	kind base.NamespaceKind
	id   base.ID
}

// Identifier implements NSContext.
func (c *nsContext) ID() base.ID {
	return c.id
}

// Kind implements NSContext.
func (c *nsContext) Kind() base.NamespaceKind {
	return c.kind
}

var _ NSContext = (*nsContext)(nil)

type internalContextKey int

const (
	nsContextKey internalContextKey = iota
)

func GetNSContext(c context.Context) NSContext {
	return c.Value(nsContextKey).(NSContext)
}

func WithDefaultNSContext(parent ctx.RequestContext, kind base.NamespaceKind, id base.ID) ctx.RequestContext {
	nsCtx := &nsContext{
		kind: kind,
		id:   id,
	}
	return parent.WithValue(nsContextKey, nsCtx)
}

var keyVaultNameIdentiferPattern = regexp.MustCompile(`^[0-9A-Za-z\-]+$`)

func VerifyKeyVaultIdentifier(id base.ID) error {
	if _, ok := id.AsUUID(); ok {
		return nil
	}
	if len(id) < 1 || len(id) > 48 || !keyVaultNameIdentiferPattern.MatchString(string(id)) {
		return fmt.Errorf("%w: invalid identifier, %s", base.ErrResponseStatusBadRequest, id)
	}

	return nil
}

func WithResovingMeNSContext(parent ctx.RequestContext, kind base.NamespaceKind, id base.ID) (ctx.RequestContext, *nsContext) {
	if strings.EqualFold("me", string(id)) {
		id = base.IDFromUUID(auth.GetAuthIdentity(parent).ClientPrincipalID())
	}
	nsCtx := &nsContext{
		kind: kind,
		id:   id,
	}
	return parent.WithValue(nsContextKey, nsCtx), nsCtx
}

func (nsCtx *nsContext) AllowSelf() authz.AuthZFunc {
	return func(c ctx.RequestContext) (ctx.RequestContext, authz.AuthzResult) {
		identity := auth.GetAuthIdentity(c)
		if base.IDFromUUID(identity.ClientPrincipalID()) == nsCtx.id {
			return c, authz.AuthzResultAllow
		}
		return c, authz.AuthzResultNone
	}
}
