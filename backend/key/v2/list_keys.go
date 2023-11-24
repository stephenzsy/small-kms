package key

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListKeys implements admin.ServerInterface.
func (*KeyAdminServer) ListKeys(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, params admin.ListKeysParams) error {
	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().
		WithExtraColumns("c.status", "c.iat", "c.exp").
		WithOrderBy("c.iat DESC")
	if params.PolicyId != nil && *params.PolicyId != "" {
		policyIdentifier := resdoc.NewDocIdentifier(
			namespaceProvider, namespaceId,
			models.ResourceProviderKeyPolicy,
			*params.PolicyId)
		qb.WithWhereClauses("c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyIdentifier.String()})
	} else {
		qb.WithExtraColumns("c.policy")
	}
	pager := resdoc.NewQueryDocPager[*KeyDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: namespaceProvider,
		NamespaceID:       namespaceId,
		ResourceProvider:  models.ResourceProviderKey,
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *KeyDoc) keymodels.KeyRef {
		return doc.ToKeyRef()
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

func ListLatestActiveKeysByPolicyInternal(c ctx.RequestContext, namespaceProvider models.NamespaceProvider, namespaceId string, policyIdentifier resdoc.DocIdentifier) ([]string, error) {
	qb := resdoc.NewDefaultCosmoQueryBuilder().
		WithExtraColumns("c.status", "c.iat", "c.exp").
		WithWhereClauses("c.status = 'active'").
		WithWhereClauses("c.policy = @policy").
		WithWhereClauses("NOT IS_DEFINED(c.exp) OR c.exp > (GetCurrentTimestamp() / 1000)").
		WithOrderBy("c.iat DESC").
		WithOffsetLimit(0, 2)

	qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyIdentifier.String()})

	pager := resdoc.NewQueryDocPager[*KeyDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: namespaceProvider,
		NamespaceID:       namespaceId,
		ResourceProvider:  models.ResourceProviderKey,
	})

	return utils.PagerToSlice(utils.NewMappedItemsPager(pager, func(doc *KeyDoc) string {
		return doc.ID
	}))

}
