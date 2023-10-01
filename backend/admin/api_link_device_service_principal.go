package admin

import (
	"context"
	"errors"
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

func (doc *NsRelDoc) toServicePrincipalLinkedDevice() (*ServicePrincipalLinkedDevice, error) {
	if doc == nil {
		return nil, nil
	}
	return &ServicePrincipalLinkedDevice{
		ApplicationClientID: DefaultIfNil(doc.Attributes.AppID),
		DeviceID:            DefaultIfNil(doc.Attributes.DeviceID),
		DeviceOID:           DefaultIfNil(doc.LinkedNamespaces.Device),
		ApplicationOID:      DefaultIfNil(doc.LinkedNamespaces.Application),
		ServicePrincipalID:  DefaultIfNil(doc.LinkedNamespaces.ServicePrincipal),
	}, nil
}

func (s *adminServer) getDeviceServicePrincipalLinkDoc(c context.Context, nsID uuid.UUID) (doc *NsRelDoc, relID uuid.UUID, err error) {
	relID = common.GetCanonicalNamespaceRelationID(nsID, common.NSRelNameDASPLink)
	doc, err = s.readNsRel(c, nsID, relID)
	return
}

func (s *adminServer) getDeviceServicePrincipalLink(c context.Context, nsID uuid.UUID) (*ServicePrincipalLinkedDevice, error) {
	relDoc, relID, err := s.getDeviceServicePrincipalLinkDoc(c, nsID)
	if err != nil {
		return nil, common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:rel:%s", nsID, relID))
	}

	return relDoc.toServicePrincipalLinkedDevice()
}

func (s *adminServer) GetDeviceServicePrincipalLinkV2(c *gin.Context, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	r, err := s.getDeviceServicePrincipalLink(c, nsID)
	if err != nil {
		if errors.Is(err, common.ErrStatusNotFound) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}

		respondInternalError(c, err, "failed to get namespace relation")
		return
	}

	c.JSON(http.StatusOK, r)
}

func (s *adminServer) createDeviceServicePrincipalLinkV2(c context.Context, nsID uuid.UUID) (*ServicePrincipalLinkedDevice, error) {
	log.Info().Msgf("createDeviceServicePrincipalLinkV2: %s - start", nsID)

	defer log.Info().Msgf("createDeviceServicePrincipalLinkV2: %s - end", nsID)

	var device msgraphmodels.Deviceable

	// device require to have a profile
	graphProfileDoc, err := s.graphService.GetGraphProfileDoc(c, nsID, kmsdoc.DocTypeExtNameDevice)
	if err != nil {
		return nil, fmt.Errorf("%w: device must be registered first", err)
	}
	log.Info().Msgf("device %s: profile loaded", nsID)

	// need to fetch device from graph
	devGraphObj, err := s.graphService.GetGraphObjectByID(c, nsID)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("%s", nsID))
		if errors.Is(err, common.ErrStatusNotFound) {
			// device is no longer available, schedule profile deletion
			if deleteErr := s.graphService.DeleteGraphProfileDoc(c, graphProfileDoc); deleteErr != nil {
				err = deleteErr
			}
		}
		return nil, err
	}
	log.Info().Msgf("device %s: loaded from msgraph", nsID)

	// verify is device
	var ok bool
	device, ok = devGraphObj.(msgraphmodels.Deviceable)
	if !ok {
		if deleteErr := s.graphService.DeleteGraphProfileDoc(c, graphProfileDoc); deleteErr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: namespace is not a device: %s", common.ErrStatusBadRequest, nsID)
	}
	// device is verified, write new object to cosmos
	deviceDoc := s.graphService.NewDeviceDocFromGraph(device)
	if err := deviceDoc.Persist(c); err != nil {
		return nil, err
	}
	log.Info().Msgf("device %s: verified and profile persisted: %s", nsID)

	// next look up existing relation
	deviceRelDoc, deviceRelID, err := s.getDeviceServicePrincipalLinkDoc(c, nsID)
	if err != nil {
		err = common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:rel:%s", nsID, deviceRelID))
		if errors.Is(err, common.ErrStatusNotFound) {
			// clear error
			deviceRelDoc = nil
			err = nil
		} else {
			return nil, err
		}
	}

	// sync device doc
	deviceDirDoc, err := s.syncDirDoc(c, nsID)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("%s", nsID))
		if errors.Is(err, common.ErrStatusNotFound) && deviceRelDoc != nil {
			// remote 404, delete relDoc
			if delErr := kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), deviceRelDoc); err != nil {
				return nil, delErr
			}
		}
		return nil, err
	}
	if deviceDirDoc.OdataType != "#microsoft.graph.device" {
		if deviceRelDoc != nil && deviceRelDoc.NamespaceID == nsID {
			if delErr := kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), deviceRelDoc); err != nil {
				return nil, delErr
			}
		}
		return nil, fmt.Errorf("%w: namespace is not a device: %s", common.ErrStatusBadRequest, nsID)
	}

	// persist start doc
	if deviceRelDoc == nil {
		deviceRelDoc = new(NsRelDoc)
		deviceRelDoc.NamespaceID = nsID
		deviceRelDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, deviceRelID)
		deviceRelDoc.Status = NsRelStatusUnknown
		deviceRelDoc.SourceNamespaceID = nsID
		deviceRelDoc.LinkedNamespaces.Device = &nsID
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), deviceRelDoc); err != nil {
			return nil, err
		}
	} else {
		deviceRelDoc.SourceNamespaceID = nsID
		deviceRelDoc.LinkedNamespaces.Device = &nsID
		if patchedDoc, err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), deviceRelDoc,
			patchNsRelDocSourceNamespaceID,
			patchNsRelDocLinkedNamespacesDevice); err != nil {
			return nil, err
		} else {
			deviceRelDoc = patchedDoc
		}
	}

	// sync created application
	appOID := DefaultIfNil(deviceRelDoc.LinkedNamespaces.Application)
	var applicationDirDoc *DirectoryObjectDoc
	if appOID == uuid.Nil {
		// no application linked, register a new one
		mApplication := msgraphmodels.NewApplication()
		mApplication.SetDisplayName(ToPtr(fmt.Sprintf("small-kms-device-%s", nsID)))
		mApplication.SetSignInAudience(ToPtr("AzureADMyOrg"))
		mApplication.SetTags([]string{fmt.Sprintf("linked-device-object-id-%s", nsID), "linked-service-small-kms"})
		mApplication.SetIsFallbackPublicClient(ToPtr(true))
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
		deviceRelDoc.LinkedNamespaces[string(NSTypeApplication)] = applicationObjectID
		if err := s.patchNsRelLinkedNamespaces(c, deviceRelDoc, string(NSTypeApplication)); err != nil {
			respondInternalError(c, err, "failed to patch relDoc")
			return
		}
		if err := s.putNsRelShadow(c, deviceRelDoc, applicationObjectID); err != nil {
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
		applicationDirDoc, err = s.syncDirDoc(c, appOID)
		if err != nil {
			if common.IsGraphODataErrorNotFound(err) || common.IsAzNotFound(err) {
				if patchErr := s.patchNsRelStatus(c, deviceRelDoc, NsRelStatusError, fmt.Sprintf("namespace is no longer available: %s", appOID)); patchErr != nil {
					log.Error().Err(patchErr).Msg("failed to patch namespace relation")
				}
				respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("linked application is no longer available: %s", appOID.String()))
				return
			}
			respondInternalError(c, err, "failed to sync application doc")
			return
		}
		if applicationDirDoc.OdataType != "#microsoft.graph.application" {
			if patchErr := s.patchNsRelStatus(c, deviceRelDoc, NsRelStatusError, fmt.Sprintf("namespace is not an application: %s", appOID)); patchErr != nil {
				log.Error().Err(patchErr).Msg("failed to patch namespace relation")
			}
			respondPublicErrorMsg(c, http.StatusBadRequest, fmt.Sprintf("namespace is not an application: %s", appOID.String()))
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
	deviceRelDoc.LinkedNamespaces[string(NSTypeServicePrincipal)] = spObjectId
	if err := s.patchNsRelLinkedNamespaces(c, deviceRelDoc, string(NSTypeServicePrincipal)); err != nil {
		respondInternalError(c, err, "failed to patch relDoc")
		return
	}
	if err := s.putNsRelShadow(c, deviceRelDoc, spObjectId); err != nil {
		respondInternalError(c, err, "failed to sync relDoc")
		return
	}

}

func (s *adminServer) CreateDeviceServicePrincipalLinkV2(c *gin.Context, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	r, err := s.createDeviceServicePrincipalLinkV2(c, nsID)
	if err != nil {
		if errors.Is(err, common.ErrStatusNotFound) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}

		respondInternalError(c, err, "failed to get namespace relation")
		return
	}

	c.JSON(http.StatusOK, r)
}
