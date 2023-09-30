package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) GetLatestCertificateByTemplateV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID, templateID uuid.UUID, params GetLatestCertificateByTemplateV2Params) {
	if !authAdminOnly(c) {
		return
	}

	certDoc, readCertDocErr := s.readCertDoc(c, nsID, kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, templateID))
	if readCertDocErr != nil {
		if common.IsAzNotFound(readCertDocErr) {
			respondPublicErrorMsg(c, http.StatusNotFound, "certificate does not exist")
			return
		}
		respondInternalError(c, readCertDocErr, "failed to load existing certificate doc")
		return
	}

	certInfo, err := s.toCertificateInfo(c, certDoc, params.IncludeCertificate, nsType, nil)
	if err != nil {
		respondInternalError(c, err, "failed to convert certificate to certificate info")
		return
	}

	c.JSON(http.StatusOK, certInfo)
}

func (s *adminServer) ListCertificatesByTemplateV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID, templateId uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	docs, err := s.listCertificateDocs(c, nsID)
	if err != nil {
		respondInternalError(c, err, fmt.Sprintf("failed to list certificates namespace: %s, template: %s", nsID, templateId))
		return
	}
	r := make([]RefWithMetadata, len(docs))
	for i, doc := range docs {
		baseDocPopulateRefWithMetadata(&doc.BaseDoc, &r[i], nsType)
		if doc.FingerprintSHA1Hex != "" {
			r[i].Metadata = map[string]string{RefPropertyKeyThumbprint: doc.FingerprintSHA1Hex}
		}
		r[i].Type = RefTypeCertificate
	}

	c.JSON(http.StatusOK, r)
}
