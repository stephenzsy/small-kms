package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

// GetCertificateTemplate implements models.ServerInterface.
func (s *server) GetCertificateTemplate(c *gin.Context, profileType models.ProfileType, profileId common.Identifier, templateId common.Identifier) {
	pc, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileId)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	res, err := s.certTemplateService.GetCertificateTemplate(pc, templateId)
	wrapResponse(c, http.StatusOK, res, err)
}

// PutCertificateTemplate implements models.ServerInterface.
func (s *server) PutCertificateTemplate(c *gin.Context,
	profileType models.ProfileType,
	profileId common.Identifier,
	templateId common.Identifier) {
	req := models.CertificateTemplateParameters{}
	err := c.BindJSON(&req)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	pc, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileId)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	res, err := s.certTemplateService.PutCertificateTemplate(pc, templateId, req)
	wrapResponse(c, http.StatusOK, res, err)
}

// ListProfiles implements models.ServerInterface.
func (s *server) ListCertificateTemplates(c *gin.Context, profileType models.ProfileType, profileId common.Identifier) {
	pc, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileId)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	res, err := s.certTemplateService.ListCertificateTemplates(pc)
	wrapResponse(c, http.StatusOK, res, err)
}

// IssueCertificateFromTemplate implements models.ServerInterface.
func (s *server) IssueCertificateFromTemplate(c *gin.Context,
	profileType models.ProfileType,
	profileId common.Identifier,
	templateId common.Identifier,
	params models.IssueCertificateFromTemplateParams) {
	ctx, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileId)
	if err != nil {
		wrapResponse[*models.CertificateInfo](c, http.StatusBadRequest, nil, err)
		return
	}
	ctx, err = s.certTemplateService.WithCertificateTemplateContext(ctx, templateId)
	if err != nil {
		wrapResponse[*models.CertificateInfo](c, http.StatusBadRequest, nil, err)
	}
	res, err := s.certService.CreateCertificateFromTemplate(ctx, params)
	wrapResponse(c, http.StatusOK, res, err)
}
