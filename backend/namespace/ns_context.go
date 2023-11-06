package ns

import (
	"context"
	"fmt"
	"regexp"

	"github.com/stephenzsy/small-kms/backend/base"
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
