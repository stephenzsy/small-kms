package profile

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

func (*ProfileServer) ListProfiles(ec echo.Context, namespaceProvider models.NamespaceProvider) error {

	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	var namespaceID string
	switch namespaceProvider {
	case models.NamespaceProviderAgent:
		namespaceID = NamespaceIDApp
	case models.NamespaceProviderRootCA,
		models.NamespaceProviderIntermediateCA:
		namespaceID = NamespaceIDCA
	case models.NamespaceProviderServicePrincipal,
		models.NamespaceProviderGroup,
		models.NamespaceProviderUser:
		namespaceID = NamespaceIDGraph
	default:
		return base.ErrResponseStatusNotFound
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().WithExtraColumns("c.displayName")
	pager := resdoc.NewQueryDocPager[*ProfileDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: models.NamespaceProviderProfile,
		NamespaceID:       namespaceID,
		ResourceProvider:  models.ResourceProvider(namespaceProvider),
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *ProfileDoc) *models.Ref {
		ref := doc.ToRef()
		return &ref
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
