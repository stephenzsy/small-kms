package agentadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	profile "github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetAgent implements ServerInterface.
func (*AgentAdminServer) GetAgent(ec echo.Context, id string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc, err := getAgentInternal(c, id)
	if err != nil {

		return err
	}
	return c.JSON(200, doc.ToProfile())
}

func getAgentInternal(c context.Context, id string) (*AgentDoc, error) {
	doc := &AgentDoc{}
	err := resdoc.GetDocService(c).Read(c, resdoc.DocIdentifier{
		PartitionKey: resdoc.PartitionKey{
			NamespaceProvider: models.NamespaceProviderProfile,
			NamespaceID:       profile.NamespaceIDApp,
			ResourceProvider:  models.ProfileResourceProviderAgent,
		},
		ID: id,
	}, doc, nil)
	if err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return doc, fmt.Errorf("%w: agent not found: %s", base.ErrResponseStatusNotFound, id)
		}
	}
	return doc, err
}
