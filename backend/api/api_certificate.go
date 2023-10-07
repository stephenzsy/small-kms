package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/models"
)

// PutCertificateTemplate implements models.ServerInterface.
func (s *server) PutCertificateTemplate(c *gin.Context,
	profileType models.ProfileType,
	profileIdentifier models.Identifier,
	templateIdentifier models.Identifier) {
	req := models.CertificateTemplateParameters{}
	err := c.BindJSON(&req)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	pc, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileIdentifier)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	res, err := s.certTemplateService.PutCertificateTemplate(pc, templateIdentifier, req)
	wrapResponse(c, http.StatusOK, res, err)
}

// ListProfiles implements models.ServerInterface.
func (s *server) ListCertificateTemplates(c *gin.Context, profileType models.ProfileType, profileId models.Identifier) {
	pc, err := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileId)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	res, err := s.certTemplateService.ListCertificateTemplates(pc)
	wrapResponse(c, http.StatusOK, res, err)
}
