package key

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListKeyPolicies implements ServerInterface.
func (*server) ListKeyPolicies(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(namespaceKind, namespaceIdentifier, base.ResourceKindKeyPolicy)
	pager := base.NewQueryDocPager[*KeyPolicyDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *KeyPolicyDoc) *KeyPolicyRef {
		r := &KeyPolicyRef{}
		d.populateModelRef(r)
		return r
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
