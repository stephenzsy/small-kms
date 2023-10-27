package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listCertPolicies(c context.Context) ([]*CertPolicyRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy)
	pager := base.NewQueryDocPager[*CertPolicyDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
