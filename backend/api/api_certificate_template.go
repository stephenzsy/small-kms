package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func (*server) RemoveKeyVaultRoleAssignment(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	templateID shared.Identifier,
	roleAssignmentID string) error {

	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	bad := func(e error) error {
		return wrapResponse[any](ctx, http.StatusNoContent, nil, e)
	}

	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return bad(err)
	}
	err = ct.DeleteKeyVaultRoleAssignment(c, roleAssignmentID)
	return wrapResponse[any](c, http.StatusNoContent, nil, err)
}

// DeleteCertificateTemplate implements models.ServerInterface.
func (*server) DeleteCertificateTemplate(ctx echo.Context, namespaceKind shared.NamespaceKind, namespaceId shared.Identifier, templateId shared.Identifier) error {
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	return wrapEchoResponse(c, errors.New("not implemented"))
}
