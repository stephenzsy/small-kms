package key

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetSecret implements ServerInterface.
func (s *server) GetKey(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, keyID base.ID) error {
	c := ec.(ctx.RequestContext)

	c, nsCtx := ns.WithResovingMeNSContext(c, nsKind, nsID)
	c, authOk := authz.Authorize(c, authz.AllowAdmin, nsCtx.AllowSelf())
	if !authOk {
		return base.ErrResponseStatusForbidden
	}

	doc := &KeyDoc{}
	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Read(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindKey, keyID), doc, nil); err != nil {
		return wrapKeyNotFoundError(err, keyID)
	}

	model := &Key{}
	doc.populateModel(model)
	return c.JSON(http.StatusOK, model)
}

func wrapKeyNotFoundError(err error, keyID base.ID) error {
	if errors.Is(err, base.ErrAzCosmosDocNotFound) {
		return fmt.Errorf("%w, key not found: %s", base.ErrResponseStatusNotFound, keyID)
	}
	return err
}
