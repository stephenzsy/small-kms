package managedapp

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
)

type ManagedAppDoc struct {
	profile.ProfileDoc
	ApplicationID        uuid.UUID `json:"applicationId"`
	ServicePrincipalID   uuid.UUID `json:"servicePrincipalId"`
	ServicePrincipalType string    `json:"servicePrincipalType,omitempty"`
}

const (
	queryColumnApplicationID      = "c.applicationId"
	queryColumnServicePrincipalID = "c.servicePrincipalId"

	patchColumnServicePrincipalID = "/servicePrincipalId"
)

const (
	namespaceIDNameManagedApp = "managed-app"
	namespaceIDNameSystemApp  = "system-app"
)

func (d *ManagedAppDoc) Init(appID uuid.UUID, displayName string, name string) {
	if d == nil {
		return
	}
	d.ProfileDoc.Init(
		base.IDFromString(name),
		base.ProfileResourceKindManagedApp,
		base.IDFromUUID(appID),
		displayName,
	)
}

func (d *ManagedAppDoc) GetAppID() uuid.UUID {
	return d.ID.UUID()
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
	if d == nil || m == nil {
		return
	}
	d.PopulateModelRef(m)
	m.ServicePrincipalType = d.ServicePrincipalType
}

var _ base.ModelRefPopulater[managedAppRefComposed] = (*ManagedAppDoc)(nil)
var _ base.ModelPopulater[ManagedApp] = (*ManagedAppDoc)(nil)
