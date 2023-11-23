package profile

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// PutProfile implements admin.ServerInterface.
func (*ProfileServer) PutProfile(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	switch namespaceProvider {
	case models.NamespaceProviderRootCA,
		models.NamespaceProviderIntermediateCA:
		// ok
	default:
		return base.ErrResponseStatusBadRequest
	}

	if err := ns.ValidateID(namespaceId); err != nil {
		return err
	}

	params := new(models.ProfileParameters)
	if err := c.Bind(params); err != nil {
		return err
	}

	doc := &ProfileDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderProfile,
				NamespaceID:       NamespaceIDCA,
				ResourceProvider:  models.ResourceProvider(namespaceProvider),
			},
			ID: namespaceId,
		},
	}
	if params.DisplayName == "" {
		doc.DisplayName = &namespaceId
	}

	resp, err := resdoc.GetDocService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
