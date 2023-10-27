package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listProfiles(c context.Context, resourceKind base.ResourceKind) ([]*ProfileRef, error) {
	ns := ns.GetNSContext(c)
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(QueryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(base.NamespaceKindProfile, ns.Identifier(), resourceKind)
	pager := base.NewQueryDocPager[*ProfileDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ProfileDoc) *ProfileRef {
		r := &ProfileRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
