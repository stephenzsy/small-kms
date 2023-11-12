package key

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// PutKeyPolicy implements ServerInterface.
func (*server) PutKeyPolicy(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, policyID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	params := &KeyPolicyParameters{}
	if err := c.Bind(params); err != nil {
		return err
	}

	c = ns.WithNSContext(c, nsKind, nsID)
	doc := &KeyPolicyDoc{}
	if err := doc.init(c, policyID, params); err != nil {
		return err
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Upsert(c, doc, nil); err != nil {
		return err
	}

	m := &KeyPolicy{}
	doc.populateModel(m)
	return c.JSON(200, m)

}
