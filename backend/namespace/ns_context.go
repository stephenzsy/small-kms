package ns

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type NSContext interface {
	Kind() base.NamespaceKind
	Identifier() base.Identifier
}

type nsContext struct {
	kind       base.NamespaceKind
	identifier base.Identifier
}

// Identifier implements NSContext.
func (c *nsContext) Identifier() base.Identifier {
	return c.identifier
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

func WithDefaultNSContext(parent ctx.RequestContext, kind base.NamespaceKind, identifier base.Identifier) ctx.RequestContext {
	nsCtx := &nsContext{
		kind:       kind,
		identifier: identifier,
	}
	return parent.WithValue(nsContextKey, nsCtx)
}
