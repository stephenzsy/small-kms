package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) GetCertificateTemplateV2(c *gin.Context, namespaceId uuid.UUID, templateId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	doc, err := s.readCertificateTemplateDoc(c, namespaceId, templateId)
	if err != nil {
		common.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc.toCertificateTemplate())
}
