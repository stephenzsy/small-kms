package cert

import (
	"net/http"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertPolicyQueryDoc struct {
	base.QueryBaseDoc
	DisplayName string `json:"displayName"`
}

// PopulateModelRef implements base.ModelRefPopulater.
func (d *CertPolicyQueryDoc) PopulateModelRef(r *CertPolicyRef) {
	if d == nil || r == nil {
		return
	}
	d.QueryBaseDoc.PopulateModelRef(&r.ResourceReference)
	r.DisplayName = d.DisplayName
}

var _ base.ModelRefPopulater[CertPolicyRef] = (*CertPolicyQueryDoc)(nil)

func apiListCertPolicies(c ctx.RequestContext) error {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy)
	pager := base.NewQueryDocPager[*CertPolicyQueryDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyQueryDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		return r
	})
	return c.JSON(http.StatusOK, utils.NewSerializableItemsPager(modelPager))
}
