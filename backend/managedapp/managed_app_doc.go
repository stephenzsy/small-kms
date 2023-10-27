package managedapp

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
)

type ManagedAppDoc struct {
	profile.ProfileDoc
	ApplicationID      uuid.UUID `json:"applicationId"`
	ServicePrincipalID uuid.UUID `json:"servicePrincipalId"`
}

const (
	queryColumnApplicationID      = "c.applicationId"
	queryColumnServicePrincipalID = "c.servicePrincipalId"

	patchColumnServicePrincipalID = "/servicePrincipalId"
)

const namespaceIDName = "managed-app"

var namespaceIdentifierManagedApp = base.StringIdentifier(namespaceIDName)

func (d *ManagedAppDoc) Init(appID uuid.UUID, displayName string) {
	if d == nil {
		return
	}
	d.ProfileDoc.Init(
		namespaceIdentifierManagedApp,
		base.ProfileResourceKindManagedApp,
		base.UUIDIdentifier(appID),
		displayName,
	)
}

func (d *ManagedAppDoc) GetAppID() uuid.UUID {
	return d.StorageID.UUID()
}

func (d *ManagedAppDoc) PopulateModelRef(m *ManagedAppRef) {
	if d == nil || m == nil {
		return
	}
	d.ProfileDoc.PopulateModelRef(&m.ProfileRef)
	m.AppID = d.GetAppID()
	m.ApplicationID = d.ApplicationID
	m.ServicePrincipalID = d.ServicePrincipalID
}

func (d *ManagedAppDoc) PopulateModel(m *ManagedApp) {
	d.PopulateModelRef(m)
}

var _ base.ModelRefPopulater[managedAppRefComposed] = (*ManagedAppDoc)(nil)
var _ base.ModelPopulater[ManagedApp] = (*ManagedAppDoc)(nil)
