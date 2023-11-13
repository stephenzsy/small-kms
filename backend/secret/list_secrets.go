package secret

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListSecrets implements ServerInterface.
func (*server) ListSecrets(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, params ListSecretsParams) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	p := api.QueryPolicyItemsParams{
		ExtraColumns: []string{api.PolicyItemsQueryColumnCreated, api.PolicyItemsQueryColumnNotAfter},
	}
	if params.PolicyId != nil {
		p.PolicyLocator = utils.ToPtr(base.NewDocLocator(nsKind, nsID, base.ResourceKindSecretPolicy,
			base.ParseID(*params.PolicyId)))
	}
	pager := api.QueryPolicyItems[*SecretDoc](
		c,
		base.NewDocNamespacePartitionKey(nsKind, nsID, base.ResourceKindSecret),
		p)

	modelPager := utils.NewMappedItemsPager(pager, func(d *SecretDoc) *SecretRef {
		r := &SecretRef{}
		d.PopulateModelRef(r)
		return r
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

func QueryLatestSecretIDIssuedByPolicy(c ctx.RequestContext, policyFullIdentifier base.DocLocator, limit uint) ([]base.ID, error) {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithOrderBy(fmt.Sprintf("%s DESC", secretDocQueryColumnCreated)).
		WithOffsetLimit(0, limit)
	qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy", "NOT IS_DEFINED(c.deleted)")
	qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyFullIdentifier.String()})
	pager := base.NewQueryDocPager[*SecretDoc](c,
		qb,
		base.NewDocNamespacePartitionKey(policyFullIdentifier.NamespaceKind(), policyFullIdentifier.NamespaceID(), base.ResourceKindSecret))

	return utils.PagerToSlice(utils.NewMappedItemsPager(pager, func(d *SecretDoc) base.ID {
		return d.ID
	}))
}
