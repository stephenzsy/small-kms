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

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(secretDocQueryColumnVersion).
		WithOrderBy(fmt.Sprintf("%s DESC", secretDocQueryColumnCreated))

	storageNsID := base.NewDocNamespacePartitionKey(nsKind, nsID, base.ResourceKindSecret)

	if params.PolicyId != nil {
		policyIdentifier := base.ParseID(*params.PolicyId)

		policyLocator := base.NewDocLocator(nsKind, nsID, base.ResourceKindSecretPolicy, policyIdentifier)

		qb = qb.WithWhereClauses("c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyLocator.String()})
	}

	pager := base.NewQueryDocPager[*SecretDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *SecretDoc) *SecretRef {
		r := &SecretRef{}
		d.PopulateModelRef(r)
		return r
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
