package agentadmin

import (
	"github.com/labstack/echo/v4"
	appadmin "github.com/stephenzsy/small-kms/backend/admin/app"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func (s *AgentAdminServer) ListAgents(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return echo.ErrForbidden
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().WithExtraColumns("c.displayName")
	pager := resdoc.NewQueryDocPager[*AgentDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: models.NamespaceProviderProfile,
		NamespaceID:       appadmin.AppNamespaceID,
		ResourceProvider:  models.ProfileResourceProviderAgent,
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *AgentDoc) *models.Ref {
		ref := doc.ToRef()
		ref.DisplayName = doc.DisplayName
		return &ref
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
