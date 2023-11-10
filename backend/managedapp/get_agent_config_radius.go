package managedapp

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetAgentConfigRadius implements ServerInterface.
func (s *server) GetAgentConfigRadius(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID) error {
	c := ec.(ctx.RequestContext)
	c, nsCtx := ns.WithResovingMeNSContext(c, nsKind, nsID)
	c, authOk := authz.Authorize(c, authz.AllowAdmin, nsCtx.AllowSelf())
	if !authOk {
		return base.ErrResponseStatusForbidden
	}

	doc, err := apiReadAgentConfigRadiusDoc(c)
	if err != nil {
		return err
	}
	m := &AgentConfigRadius{}
	doc.populateModel(m)
	return c.JSON(200, m)
}

func apiReadAgentConfigRadiusDoc(c ctx.RequestContext) (*AgentConfigRadiusDoc, error) {
	if doc, err := readAgentConfigRadiusDoc(c); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameRadius)
		}
		return nil, err
	} else {
		return doc, nil
	}
}

func readAgentConfigRadiusDoc(c ctx.RequestContext) (*AgentConfigRadiusDoc, error) {
	nsCtx := ns.GetNSContext(c)
	return readAgentConfigRadiusDocByLocator(c, base.NewDocLocator(
		nsCtx.Kind(),
		nsCtx.ID(),
		base.ResourceKindNamespaceConfig,
		base.ID(base.AgentConfigNameRadius)))
}

func readAgentConfigRadiusDocByLocator(c ctx.RequestContext, locator base.DocLocator) (*AgentConfigRadiusDoc, error) {
	doc := &AgentConfigRadiusDoc{}
	err := base.GetAzCosmosCRUDService(c).Read(c, locator, doc, nil)
	return doc, err
}
