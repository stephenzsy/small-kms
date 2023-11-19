package cert

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func (*CertServer) ListCertificatePolicies(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().WithExtraColumns(queryColumnDisplayName)
	pager := resdoc.NewQueryDocPager[*CertPolicyDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: namespaceProvider,
		NamespaceID:       namespaceId,
		ResourceProvider:  models.ResourceProviderCertPolicy,
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *CertPolicyDoc) *models.Ref {
		ref := doc.ToRef()
		ref.DisplayName = &doc.DisplayName
		return &ref
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
