package admin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
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
		certDoc.KeyInfo.Kty == models.KeyTypeRSA && (certDoc.KeyInfo.KeySize == nil || templateDoc.KeyProperties.KeySize == nil ||
			*certDoc.KeyInfo.KeySize != *templateDoc.KeyProperties.KeySize) {
		return "alg or key mismatch"
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

func (s *adminServer) IssueCertificateByTemplateV2(c *gin.Context, nsID uuid.UUID, templateID uuid.UUID, params IssueCertificateByTemplateV2Params) {
	if !authAdminOnly(c) {
		return
	}

	certDoc, readCertDocErr := s.readCertDoc(c, nsID, kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForTemplate, templateID))
	if readCertDocErr != nil {
		if !errors.Is(readCertDocErr, common.ErrStatusNotFound) {
			respondInternalError(c, readCertDocErr, "failed to load existing certificate doc")
			return
		}
		certDoc = nil
	}

	var createCertPemBlob []byte
	// create certificate
	successResponseCode := http.StatusOK
	// verify template
	templateDoc, err := s.readCertificateTemplateDoc(c, nsID, templateID)
	if err != nil {
		common.RespondError(c, err)
		return
	}
	if !templateDoc.IsActive() {
		respondPublicErrorMsg(c, http.StatusBadRequest, "template is not active")
		return
	}

	shouldApplyReason := s.shouldCreateCertificateForTemplate(c, nsID, templateDoc, certDoc)
	if len(shouldApplyReason) == 0 {
		log.Info().Msg("certificate up-to-date, no need to apply certificate template")
	} else {
		if !isAllowedCaNamespace(nsID) {
			// before issue certificate, verify both profile and object exist in graph
			gc, err := s.msGraphClient(c)
			if err != nil {
				common.RespondError(c, err)
			}
			dirObj, err := gc.DirectoryObjects().ByDirectoryObjectId(nsID.String()).Get(c, nil)
			if err != nil {
				err = common.WrapMsGraphNotFoundErr(err, "dirobj:"+nsID.String())

				common.RespondError(c, err)
				return
			}

			// store doc
			doc, err := s.graphService.GetGraphProfileDoc(c, nsID, graph.MsGraphOdataTypeAny)
			if err != nil {
				common.RespondError(c, err)
			}
			if gp, ok := dirObj.(graph.GraphProfileable); ok {
				doc = s.graphService.NewGraphProfileDoc(s.TenantID(), gp)
			}
			// update profile doc
			if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
				common.RespondError(c, err)
				return
			}
		}
		// create certificate
		certDoc, createCertPemBlob, err = s.createCertificateFromTemplate(c, nsID, &certTemplateProcessor{
			tmplDoc: templateDoc,
		})
		if err != nil {
			common.RespondError(c, err)
			return
		}
		successResponseCode = http.StatusCreated

		// psersist certificate in cosmos
		err = kmsdoc.AzCosmosCreate(c, s.AzCosmosContainerClient(), certDoc)
		if err != nil {
			respondInternalError(c, err, "failed to store certificate metadata")
			return
		}

		// persist latest certificate for template
		certDocL := *certDoc
		certDocL.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForTemplate, templateID)
		certDocL.AliasID = &certDoc.ID
		err = kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), &certDocL)
		if err != nil {
			respondInternalError(c, err, "failed to store certificate metadata for template")
			return
		}
	}

	if certDoc == nil {
		respondPublicErrorMsg(c, http.StatusNotFound, "certificate does not exist")
		return
	}

	certInfo, err := s.toCertificateInfo(c, certDoc, params.IncludeCertificate, createCertPemBlob)
	if err != nil {
		respondInternalError(c, err, "failed to convert certificate to certificate info")
		return
	}

	c.JSON(successResponseCode, certInfo)
}
