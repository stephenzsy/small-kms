package key

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListCertificates implements ServerInterface.
func (s *server) ListKeys(ec echo.Context,
	nsKind base.NamespaceKind,
	nsID base.ID, params ListKeysParams) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	p := api.QueryPolicyItemsParams{
		ExtraColumns: []string{api.PolicyItemsQueryColumnCreated, api.PolicyItemsQueryColumnNotAfter},
	}
	if params.PolicyId != nil {
		p.PolicyLocator = utils.ToPtr(base.NewDocLocator(nsKind, nsID, base.ResourceKindKeyPolicy,
			base.ParseID(*params.PolicyId)))
	}
	pager := api.QueryPolicyItems[*KeyDoc](
		c,
		base.NewDocNamespacePartitionKey(nsKind, nsID, base.ResourceKindKey),
		p)

	modelPager := utils.NewMappedItemsPager(pager, func(d *KeyDoc) *KeyRef {
		r := &KeyRef{}
		d.populateModelRef(r)
		return r
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
