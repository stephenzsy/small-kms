package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	msgraphsp "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) GetDeviceServicePrincipalLinkV2(c *gin.Context, nsID uuid.UUID, params GetDeviceServicePrincipalLinkV2Params) {
	if !authAdminOnly(c) {
		return
	}

	relDoc, err := s.readNsRel(c, nsID, common.WellKnownNSRelID_DeviceLinkServicePrincipal)
	if err != nil {
		if common.IsAzNotFound(err) {
			relDoc = nil
		} else {
			respondInternalError(c, err, "failed to get namespace relation")
			return
		}
	}

	//httpSuccessCode := http.StatusOK

	if common.ResolveBoolPtrValue(params.Apply) {

		// sync device doc
		deviceDirDoc, err := s.syncDirDoc(c, nsID)
		if err != nil {
			if common.IsGraphODataErrorNotFound(err) || common.IsAzNotFound(err) {
				if relDoc != nil && relDoc.NamespaceID == nsID {
					if patchErr := s.patchNsRelStatus(c, relDoc, NsRelStatusError, fmt.Sprintf("device not exist or not registered: %s", nsID)); patchErr != nil {
						log.Error().Err(patchErr).Msg("failed to patch namespace relation")
					}
				}
				respondPublicError(c, http.StatusNotFound, fmt.Errorf("device not exist or not registered: %s", nsID))
				return
			}
			respondInternalError(c, err, "failed to sync directory")
			return
		}

		if deviceDirDoc.OdataType != "#microsoft.graph.device" {
			if relDoc != nil && relDoc.NamespaceID == nsID {
				if patchErr := s.patchNsRelStatus(c, relDoc, NsRelStatusError, fmt.Sprintf("not a device: %s", nsID)); patchErr != nil {
					log.Error().Err(patchErr).Msg("failed to patch relDoc")
				}
			}
			respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace is not a device: %s", nsID))
			return
		}

		if relDoc == nil || relDoc.NamespaceID != nsID {
			// very rare case: a namepace registered as something else is now a device, recreate and persist doc
			relDoc = new(NsRelDoc)
			relDoc.LinkedNamespaces = make(map[string]uuid.UUID, 0)
			relDoc.NamespaceID = nsID
			relDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, common.WellKnownNSRelID_DeviceLinkServicePrincipal)
			relDoc.Status = NsRelStatusUnknown
			relDoc.SourceNamespaceID = nsID
			err := kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, relDoc)
			if err != nil {
				respondInternalError(c, err, "failed to persist relDoc")
			}
		} else if relDoc.LinkedNamespaces == nil {
			relDoc.LinkedNamespaces = make(map[string]uuid.UUID, 0)
		}

		// sync created application
		applicationID := relDoc.LinkedNamespaces[string(NSTypeApplication)]
		var applicationDirDoc *DirectoryObjectDoc
		if applicationID == uuid.Nil {
			// no application linked, register a new one
			mApplication := msgraphmodels.NewApplication()
			mApplication.SetDisplayName(ToPtr(fmt.Sprintf("small-kms-device-%s", nsID)))
			mApplication.SetSignInAudience(ToPtr("AzureADMyOrg"))
			mApplication.SetTags([]string{fmt.Sprintf("linked-device-object-id-%s", nsID), "liked-service-small-kms"})
			application, err := s.msGraphClient.Applications().Post(c, mApplication, nil)
			if err != nil {
				// failed to create application
				respondInternalError(c, err, "failed to create application")
				return
			}
			// sync application doc
			applicationIdString := *application.GetId()
			applicationObjectID, err := uuid.Parse(applicationIdString)
			if err != nil {
				respondInternalError(c, err, "failed to sync application")
				return
			}
			// patch the application object id to relDoc
			relDoc.LinkedNamespaces[string(NSTypeApplication)] = applicationObjectID
			if err := s.patchNsRelLinkedNamespaces(c, relDoc, string(NSTypeApplication)); err != nil {
				respondInternalError(c, err, "failed to patch relDoc")
				return
			}
			if err := s.putNsRelShadow(c, relDoc, applicationObjectID); err != nil {
				respondInternalError(c, err, "failed to sync relDoc")
				return
			}
			applicationDirDoc, err = s.syncDirDoc(c, applicationObjectID)
			if err != nil {
				respondInternalError(c, err, "failed to sync relDoc")
				return
			}
			// application doc synced
		} else {
			applicationDirDoc, err = s.syncDirDoc(c, applicationID)
			if err != nil {
				if common.IsGraphODataErrorNotFound(err) || common.IsAzNotFound(err) {
					if patchErr := s.patchNsRelStatus(c, relDoc, NsRelStatusError, fmt.Sprintf("namespace is no longer available: %s", applicationID)); patchErr != nil {
						log.Error().Err(patchErr).Msg("failed to patch namespace relation")
					}
					respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("linked application is no longer available: %s", applicationID.String()))
					return
				}
				respondInternalError(c, err, "failed to sync application doc")
				return
			}
			if applicationDirDoc.OdataType != "#microsoft.graph.application" {
				if patchErr := s.patchNsRelStatus(c, relDoc, NsRelStatusError, fmt.Sprintf("namespace is not an application: %s", applicationID)); patchErr != nil {
					log.Error().Err(patchErr).Msg("failed to patch namespace relation")
				}
				respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace is not an application: %s", applicationID.String()))
				return
			}
		}

		// lookup service principal
		sp, err := s.msGraphClient.ServicePrincipalsWithAppId(&applicationDirDoc.Application.AppID).Get(c, &msgraphsp.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &msgraphsp.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id"},
			},
		})
		if err != nil {
			if common.IsGraphODataErrorNotFoundWithAltErrorCode(err, 0) {
				// not found, create one
				spInput := msgraphmodels.NewServicePrincipal()
				spInput.SetAppId(&applicationDirDoc.Application.AppID)
				sp, err = s.msGraphClient.ServicePrincipals().Post(c, spInput, nil)
				if err != nil {
					respondInternalError(c, err, "failed to create service principal")
					return
				}
			} else {
				respondInternalError(c, err, "failed to get service principal")
				return
			}
		}
		spObjectId, err := uuid.Parse(*sp.GetId())
		if err != nil {
			respondInternalError(c, err, "failed to parse service principal object id")
			return
		}
		if _, err = s.syncDirDoc(c, spObjectId); err != nil {
			respondInternalError(c, err, "failed to sync service principal")
			return
		}

		// patch the application object id to relDoc
		relDoc.LinkedNamespaces[string(NSTypeServicePrincipal)] = spObjectId
		if err := s.patchNsRelLinkedNamespaces(c, relDoc, string(NSTypeServicePrincipal)); err != nil {
			respondInternalError(c, err, "failed to patch relDoc")
			return
		}
		if err := s.putNsRelShadow(c, relDoc, spObjectId); err != nil {
			respondInternalError(c, err, "failed to sync relDoc")
			return
		}

	}

	if relDoc == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no link found"})
		return
	}

	respondPublicErrorMsg(c, http.StatusNotImplemented, "not implemented")
}
