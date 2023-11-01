package api

import (
	"github.com/labstack/echo/v4"
	agentconfig "github.com/stephenzsy/small-kms/backend/agent-config"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func (*server) GetDockerInfo(ec echo.Context, namespaceId shared.Identifier) error {

	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err := ns.WithNamespaceContext(c, shared.NamespaceKindServicePrincipal, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}

	return wrapEchoResponse(c, agentconfig.ApiProxyGetDockerInfo(c))
}
