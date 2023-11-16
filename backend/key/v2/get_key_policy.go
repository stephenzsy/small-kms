package key

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetKeyPolicy implements ServerInterface.
func (*KeyAdminServer) GetKeyPolicy(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc := &KeyPolicyDoc{}
	if err := resdoc.GetDocService(c).Read(c, resdoc.DocIdentifier{
		PartitionKey: resdoc.PartitionKey{
			NamespaceProvider: namespaceProvider,
			NamespaceID:       namespaceId,
			ResourceProvider:  models.ResourceProviderKeyPolicy,
		},
		ID: id,
	}, doc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: key policy not found: %s", base.ErrResponseStatusNotFound, id)
		}
		return err
	}
	return c.JSON(200, doc.ToModel())
}
