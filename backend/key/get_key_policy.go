package key

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetKeyPolicy implements ServerInterface.
func (*server) GetKeyPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)

	doc, err := apiGetKeyPolicyDoc(c, resourceIdentifier)
	if err != nil {
		return err
	}

	m := &KeyPolicy{}
	doc.populateModel(m)
	return c.JSON(200, m)
}

func apiGetKeyPolicyDoc(c ctx.RequestContext, policyID base.ID) (*KeyPolicyDoc, error) {
	doc := &KeyPolicyDoc{}
	nsCtx := ns.GetNSContext(c)
	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Read(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindKeyPolicy, policyID), doc, nil); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: key policy %s not found", base.ErrResponseStatusNotFound, policyID)
		}
		return nil, err
	}
	return doc, nil
}
