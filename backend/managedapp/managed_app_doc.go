package managedapp

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

type ManagedAppDoc struct {
	base.BaseDoc
	DisplayName        string    `json:"displayName"`
	ApplicationID      uuid.UUID `json:"applicationId"`
	ServicePrincipalID uuid.UUID `json:"servicePrincipalId"`
}

const (
	patchColumnServicePrincipalID = "/servicePrincipalId"
)

const namespaceIDName = "managed-app"

var namespaceIdentifierManagedApp = base.StringIdentifier(namespaceIDName)

func NewManagedAppDoc(appID uuid.UUID, displayName string) *ManagedAppDoc {
	doc := &ManagedAppDoc{
		BaseDoc: base.BaseDoc{
			NamespaceKind:       base.NamespaceKindProfile,
			NamespaceIdentifier: namespaceIdentifierManagedApp,
			ResourceKind:        base.ResourceKindManagedApp,
			ResourceIdentifier:  base.UUIDIdentifier(appID),
		},
		DisplayName: displayName,
	}
	return doc
}

func (d *ManagedAppDoc) GetAppID() uuid.UUID {
	return d.ResourceIdentifier.UUID()
}

func (d *ManagedAppDoc) GetStorageID(context.Context) uuid.UUID {
	return d.GetAppID()
}

var _ base.CRUDDocHasCustomStorageID = (*ManagedAppDoc)(nil)
