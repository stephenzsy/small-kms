package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func createDefaultCertificateTemplateIDs(nsType NamespaceTypeShortName, nsID uuid.UUID) []RefWithMetadata {
	switch nsType {
	case NSTypeRootCA,
		NSTypeIntCA:
		return []RefWithMetadata{{
			ID:            uuid.Nil,
			DisplayName:   "default",
			IsActive:      utils.ToPtr(false),
			NamespaceID:   nsID,
			Type:          RefTypeCertificateTemplate,
			NamespaceType: nsType,
		}}
	case NSTypeGroup:
		spID := common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_ServicePrincipalClientCredential)
		return []RefWithMetadata{
			{
				ID:            spID,
				DisplayName:   string(common.DefaultCertTemplateName_ServicePrincipalClientCredential),
				IsActive:      utils.ToPtr(false),
				NamespaceID:   nsID,
				Type:          RefTypeCertificateTemplate,
				NamespaceType: nsType,
			}}
	}
	return nil
}

func (s *adminServer) ListCertificateTemplatesV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	docs, err := s.listCertificateTemplateDoc(c, nsID)
	if err != nil {
		respondInternalError(c, err, fmt.Sprintf("failed to list certificate templates: %s", nsID))
		return
	}
	defaultList := createDefaultCertificateTemplateIDs(nsType, nsID)
	r := make([]RefWithMetadata, len(defaultList), len(docs)+len(defaultList))
	copy(r, defaultList)
	for _, doc := range docs {
		if doc.ID.GetUUID().Version() != 4 {
			for i := range defaultList {
				if r[i].ID == doc.ID.GetUUID() {
					baseDocPopulateRefWithMetadata(&doc.BaseDoc, &r[i], nsType)
					r[i].DisplayName = doc.DisplayName
					r[i].Type = RefTypeCertificateTemplate
					r[i].IsActive = utils.ToPtr(doc.Deleted == nil || doc.Deleted.IsZero())
					goto continueLoop
				}
			}
		}
		{
			item := RefWithMetadata{}
			baseDocPopulateRefWithMetadata(&doc.BaseDoc, &item, nsType)
			item.DisplayName = doc.DisplayName
			item.Type = RefTypeCertificateTemplate
			item.IsActive = utils.ToPtr(doc.Deleted == nil || doc.Deleted.IsZero())
			r = append(r, item)
		}
	continueLoop:
	}

	c.JSON(http.StatusOK, r)
}

func (s *adminServer) GetCertificateTemplateV2(c *gin.Context, namespaceType NamespaceTypeShortName, namespaceId uuid.UUID, templateId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	doc, err := s.readCertificateTemplateDoc(c, namespaceId, templateId)
	if err != nil {
		common.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc.toCertificateTemplate(namespaceType))
}
