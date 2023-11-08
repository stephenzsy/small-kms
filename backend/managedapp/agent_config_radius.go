package managedapp

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

// GetAgentConfigRadius implements ServerInterface.
func (s *server) GetAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	panic("unimplemented")
}

// PutAgentConfigRadius implements ServerInterface.
func (s *server) PutAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	panic("unimplemented")
}

// PatchAgentConfigRadius implements ServerInterface.
func (s *server) PatchAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	panic("unimplemented")
}
