package agentadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func getAgentInstanceInternal(c context.Context, nsID, instanceID string) (*AgentInstanceDoc, error) {
	doc := &AgentInstanceDoc{}
	docSvc := resdoc.GetDocService(c)
	err := docSvc.Read(c, resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal, nsID, models.ResourceProviderAgentInstance, instanceID), doc, nil)
	if err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return doc, fmt.Errorf("%w: agent instance not found", base.ErrResponseStatusNotFound)
		}
	}
	return doc, err
}
