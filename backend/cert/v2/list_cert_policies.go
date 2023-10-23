package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listCertPolicies(c context.Context, nsKind base.NamespaceKind, namespaceIdentifier base.Identifier) ([]*CertPolicyRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder(base.ResourceKindCertPolicy).
		WithExtraColumns(queryColumnDisplayName)
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.GetDefaultStorageNamespaceID(c, nsCtx.Kind(), nsCtx.Identifier())
	pager := base.NewQueryDocPager[*CertPolicyDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		r.NID = storageNsID
		r.NamespaceKind = base.NamespaceKindProfile
		r.NamespaceIdentifier = namespaceIdentifier
		r.ResourceKind = base.ResourceKindCertPolicy
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
