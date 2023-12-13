package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// DeleteAgentInstance implements admin.ServerInterface.
func (*AgentAdminServer) DeleteAgentInstance(ec echo.Context, namespaceId string, id string, params admin.DeleteAgentInstanceParams) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Delete(c, resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal, namespaceId, models.ResourceProviderAgentInstance, id), nil)
	if err != nil {
		return err
	}
	return c.NoContent(resp.RawResponse.StatusCode)
}
