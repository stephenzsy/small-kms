package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/cert"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// GetCertificate implements models.ServerInterface.
func (s *server) GetCertificate(ec echo.Context,
	namespaceKind shared.NamespaceKind, namespaceId shared.Identifier,
	certificateId shared.Identifier,
	params models.GetCertificateParams) error {

	c := ec.(RequestContext)

	isAdmin := false
	if auth.AuthorizeAdminOnly(c) {
		isAdmin = true
	}
	namespaceId, err := ns.ResolveAuthedNamespaseID(c, namespaceKind, namespaceId)
	if err != nil && !isAdmin {
		return wrapEchoResponse(c, err)
	}

	c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}

	return wrapEchoResponse(c, cert.ApiGetCertificate(c, certificateId, params))
}

// ListCertificatesByTemplate implements models.ServerInterface.
func (s *server) ListCertificatesByTemplate(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceID shared.Identifier,
	templateID shared.Identifier) error {
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceID)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	return wrapEchoResponse(c, cert.ApiListCertificatesByTemplate(c))
}

// DeleteCertificate implements models.ServerInterface.
func (*server) DeleteCertificate(ctx echo.Context, namespaceKind shared.NamespaceKind, namespaceId shared.Identifier, certificateId shared.Identifier) error {
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}

	return wrapEchoResponse(c, cert.ApiDeleteCertificate(c, certificateId))
}
