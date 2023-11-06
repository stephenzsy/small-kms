package profile

import (
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
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

func apiListProfiles(c ctx.RequestContext, resourceKind base.ResourceKind) error {
	ns := ns.GetNSContext(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(QueryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(base.NamespaceKindProfile, ns.ID(), resourceKind)
	pager := base.NewQueryDocPager[*ProfileQueryDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ProfileQueryDoc) *ProfileRef {
		r := &ProfileRef{}
		d.PopulateModelRef(r)
		return r
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
