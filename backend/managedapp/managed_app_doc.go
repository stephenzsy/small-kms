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
	queryColumnDisplayName        = "c.displayName"
	queryColumnApplicationID      = "c.applicationId"
	queryColumnServicePrincipalID = "c.servicePrincipalId"

	patchColumnServicePrincipalID = "/servicePrincipalId"
)

const namespaceIDName = "managed-app"

var namespaceIdentifierManagedApp = base.StringIdentifier(namespaceIDName)

func getManageAppDocStorageNamespaceID(c context.Context) uuid.UUID {
	return base.GetDefaultStorageNamespaceID(c, base.NamespaceKindProfile, namespaceIdentifierManagedApp)
}

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

func managedAppDocToModel(doc *ManagedAppDoc) *ManagedApp {
	if doc == nil {
		return nil
	}
	r := new(ManagedApp)
	r.NID = doc.StorageNamespaceID
	r.RID = doc.StorageID
	r.Updated = doc.Timestamp.Time
	r.Deleted = doc.Deleted
	r.UpdatedBy = doc.UpdatedBy
	r.NamespaceKind = doc.NamespaceKind
	r.NamespaceIdentifier = doc.NamespaceIdentifier
	r.ResourceKind = doc.ResourceKind
	r.ResourceIdentifier = doc.ResourceIdentifier

	r.AppID = doc.GetAppID()
	r.ApplicationID = doc.ApplicationID
	r.DisplayName = doc.DisplayName
	r.ServicePrincipalID = doc.ServicePrincipalID
	return r
}
