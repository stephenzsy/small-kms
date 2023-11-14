package agentadmin

import (
	"github.com/labstack/echo/v4"
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
		NamespaceID:       "global",
		ResourceProvider:  models.ResourceProvider(models.ProfileResourceProviderAgent),
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *AgentDoc) *models.Ref {
		return &models.Ref{
			ID:          doc.ID,
			Updated:     doc.Timestamp.Time,
			Deleted:     doc.Deleted,
			DisplayName: doc.DisplayName,
		}
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
