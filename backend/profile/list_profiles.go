package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
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

func listProfiles(c context.Context, resourceKind base.ResourceKind) ([]*ProfileRef, error) {
	ns := ns.GetNSContext(c)
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(QueryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(base.NamespaceKindProfile, ns.Identifier(), resourceKind)
	pager := base.NewQueryDocPager[*ProfileQueryDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ProfileQueryDoc) *ProfileRef {
		r := &ProfileRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
