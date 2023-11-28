package agentadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetAgentInstance implements admin.ServerInterface.
func (*AgentAdminServer) GetAgentInstance(ec echo.Context, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}
	doc, err := getAgentInstanceInternal(c, namespaceId, id)
	if err != nil {
		return err
	}
	return c.JSON(200, doc.ToModel())
}

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
