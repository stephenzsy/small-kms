package agentpush

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

// PushAgentConfigRadius implements ServerInterface.
func (s *agentServer) PushAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, resourceId base.ID, params PushAgentConfigRadiusParams) error {
	c := ec.(ctx.RequestContext)
	s.radiusConfigManager.PullConfig()
	return c.NoContent(http.StatusNoContent)
}
