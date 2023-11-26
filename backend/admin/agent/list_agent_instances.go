package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListAgentInstances implements admin.ServerInterface.
func (*AgentAdminServer) ListAgentInstances(ec echo.Context, appID string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().WithExtraColumns("c.endpoint", "c.state", "c.configVersion", "c.buildId")
	pager := resdoc.NewQueryDocPager[*AgentInstanceDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: models.NamespaceProviderAgent,
		NamespaceID:       appID,
		ResourceProvider:  models.ResourceProviderAgentInstance,
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *AgentInstanceDoc) *agentmodels.AgentInstanceRef {
		return doc.ToModel()
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
