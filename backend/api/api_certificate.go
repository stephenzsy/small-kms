package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/cert"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetCertificateTemplate implements models.ServerInterface.
func (s *server) GetCertificateTemplate(c *gin.Context, profileType models.NamespaceKind, profileId common.Identifier, templateID common.Identifier) {
	respData, respErr := (func() (*models.CertificateTemplateComposed, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		sc := s.ServiceContext(c)
		sc, err := ns.WithNamespaceContext(sc, profileType, profileId)
		if err != nil {
			return nil, err
		}
		sc, err = ct.WithCertificateTemplateContext(sc, templateID)
		if err != nil {
			return nil, err
		}
		return ct.GetCertificateTemplate(sc)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}

// PutCertificateTemplate implements models.ServerInterface.
func (s *server) PutCertificateTemplate(c *gin.Context,
	profileType models.NamespaceKind,
	profileId common.Identifier,
	templateID common.Identifier) {
	respData, respErr := (func() (*models.CertificateTemplateComposed, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		req := models.CertificateTemplateParameters{}
		err := c.BindJSON(&req)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid input body", common.ErrStatusBadRequest)
		}

		sc := s.ServiceContext(c)
		sc, err = ns.WithNamespaceContext(sc, profileType, profileId)
		if err != nil {
			return nil, err
		}
		sc, err = ct.WithCertificateTemplateContext(sc, templateID)
		if err != nil {
			return nil, err
		}
		return ct.PutCertificateTemplate(sc, req)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}

// ListProfiles implements models.ServerInterface.
func (s *server) ListCertificateTemplates(c *gin.Context, profileType models.NamespaceKind, profileId common.Identifier) {
	respData, respErr := (func() ([]*models.CertificateTemplateRefComposed, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		sc := s.ServiceContext(c)
		sc, err := ns.WithNamespaceContext(sc, profileType, profileId)
		if err != nil {
			return nil, err
		}
		return ct.ListCertificateTemplates(sc)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}

// IssueCertificateFromTemplate implements models.ServerInterface.
func (s *server) IssueCertificateFromTemplate(c *gin.Context,
	profileType models.NamespaceKind,
	profileId common.Identifier,
	templateID common.Identifier,
	params models.IssueCertificateFromTemplateParams) {
	respData, respErr := (func() (*models.CertificateInfoComposed, error) {
		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		sc := s.ServiceContext(c)
		sc, err := ns.WithNamespaceContext(sc, profileType, profileId)
		if err != nil {
			return nil, err
		}
		sc, err = ct.WithCertificateTemplateContext(sc, templateID)
		if err != nil {
			return nil, err
		}
		return cert.IssueCertificateFromTemplate(sc, params)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}
