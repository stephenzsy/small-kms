package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/cert"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// GetCertificate implements models.ServerInterface.
func (s *server) GetCertificate(ec echo.Context, namespaceKind models.NamespaceKind, namespaceId common.Identifier, certificateId common.Identifier, params models.GetCertificateParams) error {
	bad := func(e error) error {
		return wrapResponse[*models.CertificateInfoComposed](ec, http.StatusOK, nil, e)
	}
	c := ec.(RequestContext)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}

	result, err := cert.GetCertificate(c, certificateId, params)
	return wrapResponse[*models.CertificateInfoComposed](c, http.StatusOK, result, err)
}

// IssueCertificateFromTemplate implements models.ServerInterface.
func (s *server) IssueCertificateFromTemplate(ec echo.Context,
	profileType models.NamespaceKind,
	profileId common.Identifier,
	templateID common.Identifier) error {
	c := ec.(RequestContext)
	respData, respErr := (func() (*models.CertificateInfoComposed, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		c, err := ns.WithNamespaceContext(c, profileType, profileId)
		if err != nil {
			return nil, err
		}
		c, err = ct.WithCertificateTemplateContext(c, templateID)
		if err != nil {
			return nil, err
		}
		return cert.IssueCertificateFromTemplate(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// ListCertificatesByTemplate implements models.ServerInterface.
func (s *server) ListCertificatesByTemplate(ec echo.Context,
	profileType models.NamespaceKind,
	profileId common.Identifier,
	templateID common.Identifier) error {
	c := ec.(RequestContext)
	respData, respErr := (func() ([]*shared.CertificateRef, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		c, err := ns.WithNamespaceContext(c, profileType, profileId)
		if err != nil {
			return nil, err
		}
		c, err = ct.WithCertificateTemplateContext(c, templateID)
		if err != nil {
			return nil, err
		}
		return cert.ListCertificatesByTemplate(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}
