package profile

import (
	"github.com/stephenzsy/small-kms/backend/base"
)

type ProfileQueryDoc struct {
	base.QueryBaseDoc
	DisplayName string `json:"displayName"`
}

// PopulateModelRef implements base.ModelRefPopulater.
func (d *ProfileQueryDoc) PopulateModelRef(r *ProfileRef) {
	if d == nil || r == nil {
		return
	}
	d.QueryBaseDoc.PopulateModelRef(&r.ResourceReference)
	r.DisplayName = d.DisplayName
}

var _ base.ModelRefPopulater[ProfileRef] = (*ProfileQueryDoc)(nil)
