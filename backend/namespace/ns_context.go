package ns

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

type NSContext interface {
	StorageID() uuid.UUID
	Kind() base.NamespaceKind
	Identifier() base.Identifier
}
