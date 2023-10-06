package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/models"
)

// PutCertificateTemplate implements models.ServerInterface.
func (s *server) PutCertificateTemplate(c *gin.Context,
	profileType models.ProfileType,
	profileIdentifier models.NameOrUUIDIdentifier,
	templateIdentifier models.NameOrUUIDIdentifier) {
	req := models.CertificateTemplateParameters{}
	err := c.BindJSON(&req)
	if err != nil {
		wrapResponse[*models.CertificateTemplate](c, http.StatusBadRequest, nil, err)
		return
	}
	pc := s.profileService.WithProfileContext(s.ServiceContext(c), profileType, profileIdentifier)
	res, err := s.certTemplateService.PutCertificateTemplate(pc, templateIdentifier, req)
	wrapResponse(c, http.StatusOK, res, err)
}
