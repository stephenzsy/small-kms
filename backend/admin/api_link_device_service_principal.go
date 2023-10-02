package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func (doc *NsRelDoc) toServicePrincipalLinkedDevice() *ServicePrincipalLinkedDevice {
	if doc == nil {
		return nil
	}
	return &ServicePrincipalLinkedDevice{
		ApplicationClientID: utils.NilToDefault(doc.Attributes.AppID),
		DeviceID:            utils.NilToDefault(doc.Attributes.DeviceID),
		DeviceOID:           utils.NilToDefault(doc.LinkedNamespaces.Device),
		ApplicationOID:      utils.NilToDefault(doc.LinkedNamespaces.Application),
		ServicePrincipalID:  utils.NilToDefault(doc.LinkedNamespaces.ServicePrincipal),
	}
}

func (s *adminServer) getDeviceServicePrincipalLinkDoc(c context.Context, nsID uuid.UUID) (doc *NsRelDoc, relID uuid.UUID, err error) {
	relID = common.GetCanonicalNamespaceRelationID(nsID, common.NSRelNameDASPLink)
	if doc, err = s.readNsRel(c, nsID, relID); err != nil {
		err = common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:rel:%s", nsID, relID))
	}
	return
}

func (s *adminServer) GetDeviceServicePrincipalLinkV2(c *gin.Context, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	r, _, err := s.getDeviceServicePrincipalLinkDoc(c, nsID)
	if err != nil {
		if errors.Is(err, common.ErrStatusNotFound) {
			respondPublicError(c, http.StatusNotFound, err)
			return
		}

		respondInternalError(c, err, "failed to get namespace relation")
		return
	}

	c.JSON(http.StatusOK, r.toServicePrincipalLinkedDevice())
}

func (s *adminServer) createDeviceServicePrincipalLinkDoc(c context.Context, nsID uuid.UUID) (*NsRelDoc, error) {
	log.Info().Msgf("createDeviceServicePrincipalLinkDoc: %s - start", nsID)
	defer log.Info().Msgf("createDeviceServicePrincipalLinkDoc: %s - end", nsID)

	graphClient, err := s.msGraphClient(c)
	if err != nil {
		return nil, err
	}
	c = withGraphClient(c, graphClient)

	graphAppClient, err := s.msGraphAppClient()
	if err != nil {
		return nil, err
	}

	// device require to have a profile
	graphProfileDoc, err := s.graphService.GetGraphProfileDoc(c, nsID, graph.MsGraphOdataTypeDevice)
	if err != nil {
		return nil, fmt.Errorf("%w: device must be registered first", err)
	}
	log.Info().Msgf("device %s: profile loaded", nsID)

	// need to fetch device from graph
	device, err := graphClient.Devices().ByDeviceId(nsID.String()).Get(c, nil)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("device:%s", nsID))
		if errors.Is(err, common.ErrStatusNotFound) {
			// device is no longer available, schedule profile deletion
			if deleteErr := s.graphService.DeleteGraphProfileDoc(c, graphProfileDoc); deleteErr != nil {
				err = deleteErr
			}
		}
		return nil, err
	}
	log.Info().Msgf("device %s: loaded from msgraph", nsID)

	// device is verified, write new object to cosmos
	deviceDoc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), device, graph.MsGraphOdataTypeDevice)
	if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), deviceDoc); err != nil {
		return nil, err
	}

	log.Info().Msgf("device %s: verified and profile persisted", nsID)

	// next look up existing relation
	relDoc, deviceRelID, err := s.getDeviceServicePrincipalLinkDoc(c, nsID)
	if err != nil {
		err = common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:rel:%s", nsID, deviceRelID))
		if errors.Is(err, common.ErrStatusNotFound) {
			// clear error
			relDoc = nil
			err = nil
			log.Info().Msgf("device link %s: not existing", deviceRelID)
		} else {
			return nil, err
		}
	} else {
		log.Info().Msgf("device link %s: existing loaded", deviceRelID)
	}

	// patch relations docs
	deviceDeviceID, err := uuid.Parse(utils.NilToDefault(device.GetDeviceId()))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse deviceId", err)
	}
	if relDoc == nil {
		relDoc = new(NsRelDoc)
		relDoc.NamespaceID = nsID
		relDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, deviceRelID)
		relDoc.Status = NsRelStatusPending
		relDoc.SourceNamespaceID = nsID
		relDoc.LinkedNamespaces.Device = &nsID
		relDoc.Attributes.DeviceID = &deviceDeviceID
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), relDoc); err != nil {
			return nil, err
		}
	} else {
		relDoc.SourceNamespaceID = nsID
		relDoc.LinkedNamespaces.Device = &nsID
		relDoc.Attributes.DeviceID = &deviceDeviceID
		if err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), relDoc,
			patchNsRelDocSourceNamespaceID,
			patchNsRelDocLinkedNamespacesDevice); err != nil {
			return nil, err
		}
	}
	log.Info().Msgf("device link %s: patched device: %s", deviceRelID, nsID)

	// look up application
	var appObj msgraphmodels.Applicationable
	applicationID := utils.NilToDefault(relDoc.LinkedNamespaces.Application)
	if applicationID != uuid.Nil {
		if appObj, err = graphClient.Applications().ByApplicationId(applicationID.String()).Get(c, nil); err != nil {
			err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("application:%s", applicationID))
			if !errors.Is(err, common.ErrStatusNotFound) {
				return nil, err
			}
			// not found, let appObj continue to be nil
			log.Info().Msgf("device link %s: application not exist: %s", deviceRelID, applicationID)
		} else {
			log.Info().Msgf("device link %s: application loaded: %s", deviceRelID, applicationID)
		}
	}

	if appObj == nil {
		// create new
		mApplication := msgraphmodels.NewApplication()
		mApplication.SetDisplayName(ToPtr(fmt.Sprintf("small-kms-device-%s", nsID)))
		mApplication.SetSignInAudience(ToPtr("AzureADMyOrg"))
		mApplication.SetTags([]string{fmt.Sprintf("linked-device-object-id-%s", nsID), "linked-service-small-kms"})
		mApplication.SetIsFallbackPublicClient(ToPtr(true))
		if appObj, err = graphAppClient.Applications().Post(c, mApplication, nil); err != nil {
			return nil, err
		}
		if applicationID, err = uuid.Parse(utils.NilToDefault(appObj.GetId())); err != nil {
			return nil, fmt.Errorf("%w: failed to parse application id: %s", err, applicationID)
		}
		log.Info().Msgf("device link %s: application created: %s", deviceRelID, applicationID)
	}
	applicationAppID, err := uuid.Parse(utils.NilToDefault(appObj.GetAppId()))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse appId: %s", err, applicationAppID)
	}
	relDoc.LinkedNamespaces.Application = &applicationID
	relDoc.Attributes.AppID = &applicationAppID
	if err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), relDoc,
		patchNsRelDocLinkedNamespacesApplication); err != nil {
		return nil, err
	}
	log.Info().Msgf("device link %s: patched application: %s", deviceRelID, applicationID)

	// create a link dock for application
	appLinkDoc := new(NsRelDoc)
	appLinkDoc.NamespaceID = applicationID
	appLinkDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, common.GetCanonicalNamespaceRelationID(applicationID, common.NSRelNameDASPLink))
	appLinkDoc.SourceNamespaceID = nsID
	appLinkDoc.Status = NsRelStatusLink
	if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), appLinkDoc); err != nil {
		return nil, err
	}
	log.Info().Msgf("device link %s: application link created: %s", deviceRelID, applicationID)

	// look up service principal
	var spObj msgraphmodels.ServicePrincipalable

	if spObj, err = graphClient.ServicePrincipalsWithAppId(ToPtr(applicationAppID.String())).Get(c, nil); err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("servicePrincipal-appid:%s", applicationAppID))
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, err
		}
		log.Info().Msgf("device link %s: service principal not exist for app id: %s", deviceRelID, applicationAppID)

	}

	if spObj == nil {
		// create new
		mSp := msgraphmodels.NewServicePrincipal()
		mSp.SetAppId(appObj.GetAppId())

		if spObj, err = graphAppClient.ServicePrincipals().Post(c, mSp, nil); err != nil {
			return nil, err
		}
		log.Info().Msgf("device link %s: service principal created for appId: %s", deviceRelID, applicationAppID)
	}
	spID, err := uuid.Parse(utils.NilToDefault(spObj.GetId()))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse service principal id: %s", err, spID)
	}

	relDoc.LinkedNamespaces.ServicePrincipal = &spID
	relDoc.Status = NsRelStatusEnabled
	if err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), relDoc,
		patchNsRelDocLinkedNamespacesServicePrincipal); err != nil {
		return nil, err
	}
	log.Info().Msgf("device link %s: patched service principal: %s", deviceRelID, spID)

	// create a link dock for application
	spLinkDoc := new(NsRelDoc)
	spLinkDoc.NamespaceID = spID
	spLinkDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, common.GetCanonicalNamespaceRelationID(spID, common.NSRelNameDASPLink))
	spLinkDoc.SourceNamespaceID = nsID
	spLinkDoc.Status = NsRelStatusLink
	if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), appLinkDoc); err != nil {
		return nil, err
	}
	log.Info().Msgf("device link %s: service principal link created: %s", deviceRelID, spID)

	return relDoc, nil
}

func (s *adminServer) CreateDeviceServicePrincipalLinkV2(c *gin.Context, nsID uuid.UUID) {
	if !authAdminOnly(c) {
		return
	}

	r, err := s.createDeviceServicePrincipalLinkDoc(c, nsID)
	if err != nil {
		common.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, r)
}
