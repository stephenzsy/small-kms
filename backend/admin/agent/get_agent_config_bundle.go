package agentadmin

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// GetAgentConfigBundle implements admin.ServerInterface.
func (*AgentAdminServer) GetAgentConfigBundle(ec echo.Context, namespaceId string) error {
	c := ec.(ctx.RequestContext)

	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	docSvc := resdoc.GetDocService(c)
	docIdentifier := bundleDocIdentifier(namespaceId)
	doc := &AgentConfigBundleDoc{}
	if err := docSvc.Read(c, docIdentifier, doc, nil); err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: agent config bundle not found", base.ErrResponseStatusNotFound)
		}
		return err
	}

	return c.JSON(200, doc.ToModel())
}
