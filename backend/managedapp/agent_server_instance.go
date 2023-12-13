package managedapp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentInstanceDoc struct {
	base.BaseDoc
	AgentInstanceFields
}

func (d *AgentInstanceDoc) init(nsKind base.NamespaceKind, nsID base.ID, rID base.ID, req AgentInstanceFields) {
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindAgentInstance, rID)
	d.AgentInstanceFields = req
}

func (d *AgentInstanceDoc) PopulateModel(r *AgentInstance) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&r.ResourceReference)
	r.AgentInstanceFields = d.AgentInstanceFields
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
	pager := base.NewQueryDocPager[*AgentInstanceQueryDoc](c, qb, base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindAgentInstance))
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(pager))
}

// wraps not found with 404
func ApiReadAgentInstanceDoc(c ctx.RequestContext, instanceID base.ID) (*AgentInstanceDoc, error) {
	nsCtx := ns.GetNSContext(c)
	doc := &AgentInstanceDoc{}
	docSvc := base.GetAzCosmosCRUDService(c)
	err := docSvc.Read(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindAgentInstance, instanceID), doc, nil)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: agent instance with id %s not found", base.ErrResponseStatusNotFound, instanceID)
		}
		return nil, err
	}
	return doc, err
}

func apiGetAgentInstance(c ctx.RequestContext, instanceID base.ID) error {
	doc, err := ApiReadAgentInstanceDoc(c, instanceID)
	if err != nil {
		return err
	}
	m := &AgentInstance{}
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}
