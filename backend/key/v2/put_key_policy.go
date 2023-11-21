package key

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func (*KeyAdminServer) PutKeyPolicy(ec echo.Context, nsProvider models.NamespaceProvider, nsID string, ID string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	req := new(keymodels.CreateKeyPolicyRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := ns.ValidateID(ID); err != nil {
		return err
	}

	doc := &KeyPolicyDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: nsProvider,
				NamespaceID:       nsID,
				ResourceProvider:  models.ResourceProviderKeyPolicy,
			},
			ID: ID,
		},
	}
	if err := doc.init(c, req); err != nil {
		return err
	}

	resp, err := resdoc.GetDocService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
