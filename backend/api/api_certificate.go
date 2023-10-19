package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/cert"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
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
	if auth.AuthorizeAdminOnly(c) == nil {
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

// IssueCertificateFromTemplate implements models.ServerInterface.
func (s *server) IssueCertificateFromTemplate(ctx echo.Context,
	profileType shared.NamespaceKind,
	profileId shared.Identifier,
	templateID shared.Identifier,
	params models.IssueCertificateFromTemplateParams) error {
	bad := func(e error) error {
		return wrapResponse[*shared.CertificateInfo](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)
	}

	c, err := ns.WithNamespaceContext(c, profileType, profileId)
	if err != nil {
		return bad(err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return bad(err)
	}
	resp, err := cert.IssueCertificateFromTemplate(c, params)
	return wrapResponse(c, http.StatusOK, resp, err)
}

// ListCertificatesByTemplate implements models.ServerInterface.
func (s *server) ListCertificatesByTemplate(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceID shared.Identifier,
	templateID shared.Identifier) error {
	c := ctx.(RequestContext)
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return wrapEchoResponse(c, err)
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

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return wrapEchoResponse(c, err)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}

	return wrapEchoResponse(c, cert.ApiDeleteCertificate(c, certificateId))
}
