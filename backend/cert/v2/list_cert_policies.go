package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listCertPolicies(c context.Context) ([]*CertPolicyRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder(base.ResourceKindCertPolicy).
		WithExtraColumns(queryColumnDisplayName)
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.GetDefaultStorageNamespaceID(c, nsCtx.Kind(), nsCtx.Identifier())
	pager := base.NewQueryDocPager[*CertPolicyDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		r.Id.NID = storageNsID
		r.NamespaceKind = nsCtx.Kind()
		r.NamespaceIdentifier = nsCtx.Identifier()
		r.ResourceKind = base.ResourceKindCertPolicy
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
