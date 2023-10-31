package managedapp

import (
	"net/http"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentInstanceDoc struct {
	base.BaseDoc
	AgentInstanceFields
}

func (d *AgentInstanceDoc) init(nsKind base.NamespaceKind, nsID base.Identifier, rID base.Identifier, req AgentInstanceFields) {
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindAgentInstance, rID)
	d.AgentInstanceFields = req
}

func apiPutAgentInstance(c ctx.RequestContext, instanceID base.Identifier, req AgentInstanceFields) error {
	nsCtx := ns.GetNSContext(c)
	doc := &AgentInstanceDoc{}
	doc.init(nsCtx.Kind(), nsCtx.Identifier(), instanceID, req)

	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Upsert(c, doc, nil); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

type AgentInstanceQueryDoc struct {
	base.QueryBaseDoc
	AgentInstanceFields
}

func apiListAgentInstances(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns("c.endpoint", "c.version", "c.buildId").
		WithOrderBy("c.ts DESC")
	docSvc := base.GetAzCosmosCRUDService(c)
	pager := base.NewQueryDocPager[*AgentInstanceQueryDoc](docSvc, qb, base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindAgentInstance))
	sPager := utils.NewSerializableItemsPager(c, pager)
	return c.JSON(http.StatusOK, sPager)
}
