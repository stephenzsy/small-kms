package profile

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// SyncProfile implements admin.ServerInterface.
func (*ProfileServer) GetProfile(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string) error {

	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	doc := &ProfileDoc{}
	err := resdoc.GetDocService(c).Read(c, resdoc.NewDocIdentifier(models.NamespaceProviderProfile, NamespaceIDGraph, models.ResourceProvider(namespaceProvider), namespaceId), doc, nil)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return base.ErrResponseStatusNotFound
		}
		return err
	}
	return c.JSON(http.StatusOK, doc.ToModel())
}
