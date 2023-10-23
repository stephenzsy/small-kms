package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

type ProfileDoc struct {
	base.BaseDoc
	DisplayName string `json:"displayName"`
}

const (
	QueryColumnDisplayName = "c.displayName"
)

func getProfileDocStorageNamespaceID(c context.Context, namespaceIdentifier base.Identifier) uuid.UUID {
	return base.GetDefaultStorageNamespaceID(c, base.NamespaceKindProfile, namespaceIdentifier)
}

func (d *ProfileDoc) Init(
	nsID base.Identifier,
	rKind base.ResourceKind,
	rID base.Identifier,
	displayName string) {
	if d == nil {
		return
	}
	d.NamespaceKind = base.NamespaceKindProfile
	d.NamespaceIdentifier = nsID
	d.ResourceKind = rKind
	d.ResourceIdentifier = rID
	d.DisplayName = displayName
}

func (d *ProfileDoc) PopulateModelRef(m *ProfileRef) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.DisplayName = d.DisplayName
}

func (d *ProfileDoc) PopulateModel(m *Profile) {
	d.PopulateModelRef(m)
}

var _ base.ModelRefPopulater[ProfileRef] = (*ProfileDoc)(nil)
var _ base.ModelPopulater[Profile] = (*ProfileDoc)(nil)
