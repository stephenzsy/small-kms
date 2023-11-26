package key

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetKey implements admin.ServerInterface.
func (*KeyAdminServer) GetKey(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string, params admin.GetKeyParams) error {
	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	doc, err := GetKeyInternal(c, namespaceProvider, namespaceId, id)
	if err != nil {
		return err
	}

	includeJwk := false
	if params.IncludeJwk != nil {
		includeJwk = *params.IncludeJwk
	}
	model := doc.ToModel(includeJwk)
	return c.JSON(http.StatusOK, model)
}

func GetKeyInternal(c context.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) (*KeyDoc, error) {
	certDoc := &KeyDoc{}
	if err := resdoc.GetDocService(c).Read(c, resdoc.NewDocIdentifier(namespaceProvider, namespaceId, models.ResourceProviderKey, id), certDoc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: key ID: %s", base.ErrResponseStatusNotFound, id)
		}
		return nil, err
	}
	return certDoc, nil
}
