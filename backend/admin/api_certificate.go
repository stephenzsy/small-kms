package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) shouldCreateCertificateForTemplate(ctx context.Context, nsID uuid.UUID, templateDoc *CertificateTemplateDoc, certDoc *CertDoc) (renewReason string) {
	// load existing certificate
	if !certDoc.IsActive() {
		return "existing certificate does not exist or is not active"
	}

	// verify template matches certificate metadata
	if certDoc.TemplateID != templateDoc.ID {
		return "template mismatch"
	}
	if certDoc.IssuerNamespaceID != templateDoc.IssuerNamespaceID {
		return "issuer namespace mismatch"
	}
	if certDoc.SubjectBase != templateDoc.Subject.String() {
		return "subject mismatch"
	}
	if certDoc.KeyInfo.Alg == nil || *certDoc.KeyInfo.Alg != templateDoc.KeyProperties.Alg ||
		certDoc.KeyInfo.Kty == KeyTypeRSA && (certDoc.KeyInfo.KeySize == nil || templateDoc.KeyProperties.KeySize == nil ||
			*certDoc.KeyInfo.KeySize != *templateDoc.KeyProperties.KeySize) {
		return "alg or key mismatch"
	}
	if !certDoc.SubjectAlternativeNames.Equals(templateDoc.SubjectAlternativeNames) {
		return "subject alternative names mismatch"
	}
	if certDoc.Usage != templateDoc.Usage {
		return "usage mismatch"
	}

	// verify life time
	if templateDoc.LifetimeTrigger.DaysBeforeExpiry != nil {
		daysBeforeExpiry := *templateDoc.LifetimeTrigger.DaysBeforeExpiry
		if daysBeforeExpiry > 0 && time.Now().AddDate(0, 0, int(daysBeforeExpiry)).
			After(certDoc.NotAfter) {
			return "within days before expiry"
		}
	} else if templateDoc.LifetimeTrigger.LifetimePercentage != nil {
		p := *templateDoc.LifetimeTrigger.LifetimePercentage
		if time.Now().
			After(certDoc.NotBefore.
				Add(certDoc.NotAfter.Sub(certDoc.NotBefore) * time.Duration(p) / 100)) {
			return "outside lifetime percentage"
		}
	}
	return
}

func (s *adminServer) GetCertificateV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID, templateID uuid.UUID, certID uuid.UUID, params GetCertificateV2Params) {
	if !authAdminOnly(c) {
		return
	}

	var certDocID kmsdoc.KmsDocID
	if certID == uuid.Nil {
		// use template ID
		certDocID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, templateID)
	} else {
		certDocID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, certID)
	}
	certDoc, readCertDocErr := s.readCertDoc(c, nsID, certDocID)
	if readCertDocErr != nil {
		if !common.IsAzNotFound(readCertDocErr) {
			respondInternalError(c, readCertDocErr, "failed to load existing certificate doc")
			return
		}
		certDoc = nil
	}

	var createCertPemBlob []byte
	// create certificate
	if params.Apply != nil && *params.Apply {
		// verify template
		templateDoc, err := s.readCertificateTemplateDoc(c, nsID, templateID)
		if err != nil {
			if common.IsAzNotFound(err) {
				respondPublicError(c, http.StatusNotFound, err)
				return
			}
			respondInternalError(c, err, "failed to load template doc")
			return
		}
		if !templateDoc.IsActive() {
			respondPublicErrorMsg(c, http.StatusBadRequest, "template is not active")
			return
		}

		if certID != uuid.Nil && certDoc != nil {
			respondPublicErrorMsg(c, http.StatusConflict, "certificate already exists, cannot apply certificate template")
		}

		shouldApplyReason := s.shouldCreateCertificateForTemplate(c, nsID, templateDoc, certDoc)
		if len(shouldApplyReason) == 0 {
			log.Info().Msg("certificate up-to-date, no need to apply certificate template")
		} else {

			// create certificate
			certDoc, createCertPemBlob, err = s.createCertificateFromTemplate(c, nsType, nsID, templateDoc, certID)
			if err != nil {
				if common.IsAzNotFound(err) {
					respondPublicError(c, http.StatusNotFound, err)
					return
				}
				respondInternalError(c, err, "failed to create certificate from template")
				return
			}

			// psersist certificate in cosmos
			err = kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, certDoc)
			if err != nil {
				respondInternalError(c, err, "failed to store certificate metadata")
				return
			}

			// persist latest certificate for template
			certDocL := *certDoc
			certDocL.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, templateID)
			certDocL.AliasID = &certDoc.ID
			err = kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, &certDocL)
			if err != nil {
				respondInternalError(c, err, "failed to store certificate metadata for template")
				return
			}
		}
	}
	if certDoc == nil {
		respondPublicErrorMsg(c, http.StatusNotFound, "certificate does not exist")
		return
	}

	certInfo, err := s.toCertificateInfo(c, certDoc, params.IncludeCertificate, nsType, createCertPemBlob)
	if err != nil {
		respondInternalError(c, err, "failed to convert certificate to certificate info")
		return
	}

	c.JSON(http.StatusOK, certInfo)
}

func (s *adminServer) ListCertificatesV2(c *gin.Context, nsType NamespaceTypeShortName, nsID uuid.UUID, templateId uuid.UUID) {
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
