package managedapp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type AgentInstanceDoc struct {
	base.BaseDoc
	AgentInstanceFields
}

func (d *AgentInstanceDoc) init(nsKind base.NamespaceKind, nsID base.Identifier, s string) {
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindAgentInstance, base.StringIdentifier(s))
}

func apiPutAgentConfigServerInstance(c ctx.RequestContext, instanceID string, req AgentInstanceFields) error {

	if len(instanceID) > 8 {
		return fmt.Errorf("%w: instance ID too long", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := &AgentInstanceDoc{}
	doc.init(nsCtx.Kind(), nsCtx.Identifier(), "instance-"+instanceID)
	doc.AgentInstanceFields = req

	if doc.Hostname == "" || strings.EqualFold(doc.Hostname, "localhost") {
		doc.Hostname = c.RealIP()
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	docSvc.Upsert(c, doc, nil)

	return c.NoContent(http.StatusCreated)
}
