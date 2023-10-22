package base

import (
	"context"

	"github.com/google/uuid"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type ContextKey int

const (
	SiteUrlContextKey ContextKey = iota
	AzCosmosCRUDDocServiceContextKey
	StorageNamespaceIDFunc
)

func WithStorageNamespaceIDFunc(
	parent ctx.RequestContext,
	f *func(context.Context, NamespaceKind, Identifier) *uuid.UUID) ctx.RequestContext {
	return parent.WithValue(StorageNamespaceIDFunc, f)
}

func GetStorageNamespaceIDFunc(c context.Context) *func(context.Context, NamespaceKind, Identifier) *uuid.UUID {
	if f, ok := c.Value(StorageNamespaceIDFunc).(*func(context.Context, NamespaceKind, Identifier) *uuid.UUID); ok {
		return f
	}
	return nil
}
