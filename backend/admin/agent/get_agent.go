package agentadmin

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	appadmin "github.com/stephenzsy/small-kms/backend/admin/app"
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

	doc := &AgentDoc{}
	if err := resdoc.GetDocService(c).Read(c, resdoc.DocIdentifier{
		PartitionKey: resdoc.PartitionKey{
			NamespaceProvider: models.NamespaceProviderProfile,
			NamespaceID:       appadmin.AppNamespaceID,
			ResourceProvider:  models.ProfileResourceProviderAgent,
		},
		ID: id,
	}, doc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: agent not found: %s", base.ErrResponseStatusNotFound, id)
		}
		return err
	}
	return c.JSON(200, doc.ToApplicationByAppId())
}
