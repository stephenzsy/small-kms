package cert

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListCertPolicies implements ServerInterface.
func (s *server) ListCertPolicies(ec echo.Context, nsKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	c, nsCtx, err := s.allowGeneralNonAdminAuth(c, nsKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy)
	pager := base.NewQueryDocPager[*CertPolicyQueryDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyQueryDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		return r
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

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
