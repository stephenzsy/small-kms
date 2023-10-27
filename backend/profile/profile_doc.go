package profile

import (
	"github.com/stephenzsy/small-kms/backend/base"
)

type ProfileDoc struct {
	base.BaseDoc
	DisplayName string `json:"displayName"`
}

type ProfileCRUDDoc interface {
	base.BaseDocument
	base.ModelPopulater[Profile]
	base.ModelRefPopulater[ProfileRef]
}

const (
	QueryColumnDisplayName = "c.displayName"
)

func (d *ProfileDoc) Init(
	nsID base.Identifier,
	rKind base.ResourceKind,
	rID base.Identifier,
	displayName string) {
	if d == nil {
		return
	}
	d.BaseDoc.Init(base.NamespaceKindProfile, nsID, rKind, rID)
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

var _ ProfileCRUDDoc = (*ProfileDoc)(nil)
var _ ProfileCRUDDoc = (*ProfileDoc)(nil)
