package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ListCertificateTemplatesV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	docs, err := s.listCertificateTemplateDoc(c, nsID)
	if err != nil {
		respondInternalError(c, err, fmt.Sprintf("failed to list certificate templates: %s", nsID))
		return
	}
	r := make([]Ref, len(docs))
	for i, doc := range docs {
		baseDocPopulateRef(&doc.BaseDoc, &r[i], nsType)
		r[i].DisplayName = doc.DisplayName
		r[i].Type = RefTypeCertificateTemplate
	}

	c.JSON(http.StatusOK, r)
}

func (s *adminServer) GetCertificateTemplateV2(c *gin.Context, namespaceType NamespaceTypeShortName, namespaceId uuid.UUID, templateId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	doc, err := s.readCertificateTemplateDoc(c, namespaceId, templateId)
	if err != nil {
		if common.IsAzNotFound(err) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}
		respondInternalError(c, err, fmt.Sprintf("failed to get certificate template: %s", templateId))
		return
	}

	c.JSON(http.StatusOK, doc.toCertificateTemplate(namespaceType))
}
