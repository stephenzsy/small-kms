package admin

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphapplications "github.com/microsoftgraph/msgraph-sdk-go/applications"
	msgraphdevices "github.com/microsoftgraph/msgraph-sdk-go/devices"
	msgraphdirectoryobjects "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	msgraphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	msgraphsp "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// extract owner name

func (s *adminServer) verifyDevice(c context.Context, objectID uuid.UUID, params *map[TemplateVarName]string) (msgraphmodels.Deviceable, map[TemplateVarName]string, error) {
	obj, err := s.MsGraphClient().Devices().ByDeviceId(objectID.String()).Get(c, &msgraphdevices.DeviceItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphdevices.DeviceItemRequestBuilderGetQueryParameters{
			Select: graph.GetProfileGraphSelectDeviceDoc(),
		},
	})
	if err != nil {
		return obj, nil, err
	}
	return obj, map[TemplateVarName]string{
		TemplateVarNameDeviceURI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/devices/%s", *obj.GetId()),
		TemplateVarNameDeviceAltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/devices(deviceId='{%s}')", *obj.GetDeviceId()),
	}, nil
}

func (s *adminServer) verifyApplication(c context.Context, objectID uuid.UUID, params *map[TemplateVarName]string) (msgraphmodels.Applicationable, map[TemplateVarName]string, error) {
	obj, err := s.MsGraphClient().Applications().ByApplicationId(objectID.String()).Get(c, &msgraphapplications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphapplications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, nil, err
	}
	return obj, map[TemplateVarName]string{
		TemplateVarNameApplicationURI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/applications/%s", *obj.GetId()),
		TemplateVarNameApplicationAltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/applications(appId='{%s}')", *obj.GetAppId()),
	}, nil
}

func (s *adminServer) verifyServicePrincipal(c context.Context, objectID uuid.UUID, params *map[TemplateVarName]string) (msgraphmodels.ServicePrincipalable, map[TemplateVarName]string, error) {
	obj, err := s.MsGraphClient().ServicePrincipals().ByServicePrincipalId(objectID.String()).Get(c, &msgraphsp.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphsp.ServicePrincipalItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, nil, err
	}
	return obj, map[TemplateVarName]string{
		TemplateVarNameServicePrincipalURI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals/%s", *obj.GetId()),
		TemplateVarNameServicePrincipalAltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals(appId='{%s}')", *obj.GetAppId()),
	}, nil
}

func (s *adminServer) verifyGroup(c context.Context, objectID uuid.UUID, params *map[TemplateVarName]string) (msgraphmodels.Groupable, map[TemplateVarName]string, error) {
	obj, err := s.MsGraphClient().Groups().ByGroupId(objectID.String()).Get(c, &msgraphgroups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphgroups.GroupItemRequestBuilderGetQueryParameters{
			Select: graph.GetProfileGraphSelectGroupDoc(),
		},
	})
	if err != nil {
		return obj, nil, err
	}
	return obj, map[TemplateVarName]string{
		TemplateVarNameGroupURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s", *obj.GetId()),
	}, nil
}

func (s *adminServer) verifyGroupMembership(c context.Context, objectID uuid.UUID, groupID uuid.UUID, params *map[TemplateVarName]string) (bool, error) {
	requestBody := msgraphdirectoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	requestBody.SetGroupIds([]string{groupID.String()})

	checkMemberGroups, err := s.MsGraphClient().DirectoryObjects().ByDirectoryObjectId(objectID.String()).CheckMemberGroups().Post(c, requestBody, nil)
	if err != nil {
		return false, err
	}
	memberGroups := checkMemberGroups.GetValue()
	for _, memberGroup := range memberGroups {
		if id, err := uuid.Parse(memberGroup); err == nil && id == groupID {
			return true, nil
		}
	}
	return false, nil
}

func (s *adminServer) processBeginEnrollCertForDASPLink(c context.Context, nsID uuid.UUID, templateId uuid.UUID, req CertificateEnrollmentRequestDeviceLinkedServicePrincipal) error {
	log.Info().Msgf("enroll cert for dasp link - begin: %s", req.DeviceLinkID)
	defer log.Info().Msgf("enroll cert for dasp link - end: %s", req.DeviceLinkID)

	// first check if AppID match
	authCtx, ok := auth.GetAuthIdentity(c)
	if !ok {
		return fmt.Errorf("%w: no auth identity", common.ErrStatusUnauthorized)
	}
	if req.AppID == uuid.Nil {
		return fmt.Errorf("%w: application id is nil", common.ErrStatusBadRequest)
	}
	claimedAppId := authCtx.AppIDClaim()
	if claimedAppId != req.AppID {
		return fmt.Errorf("%w: appid is required in the authorization claim", common.ErrStatusUnauthorized)
	}

	// look up relDoc
	if req.DeviceLinkID != common.GetCanonicalNamespaceRelationID(req.DeviceNamespaceID, common.NSRelNameDASPLink) {
		return fmt.Errorf("%w: device link id is invalid", common.ErrStatusBadRequest)
	}
	relDoc, err := s.readNsRel(c, req.DeviceNamespaceID, req.DeviceLinkID)
	if err != nil {
		return err
	}
	if relDoc.Status != NsRelStatusEnabled {
		return fmt.Errorf("%w: device link is not enabled", common.ErrStatusBadRequest)
	}
	if deviceID, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.Device); !nonNil || deviceID != req.DeviceNamespaceID {
		return fmt.Errorf("%w: device object id invalid", common.ErrStatusBadRequest)
	}
	if spID, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.ServicePrincipal); !nonNil || spID != req.ServicePrincipalID {
		return fmt.Errorf("%w: service principal id dinvalid", common.ErrStatusBadRequest)
	}
	if appID, nonNil := utils.NonNilUUID(relDoc.Attributes.AppID); !nonNil || appID != req.AppID {
		return fmt.Errorf("%w: application client id does invalid", common.ErrStatusBadRequest)
	}
	if _, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.Application); !nonNil {
		return fmt.Errorf("%w: application object id is nil", common.ErrStatusBadRequest)
	}
	log.Info().Msgf("link doc loaded and verified: %s", req.DeviceLinkID)

	// prep parameters
	params := make(map[TemplateVarName]string)
	// verify against ms graph
	if obj, deviceParams, err := s.verifyDevice(c, req.DeviceNamespaceID, &params); err != nil {
		return err
	} else {
		maps.Copy(params, deviceParams)
		doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeDevice)
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("device verified: %s", req.DeviceNamespaceID)

	if _, appParams, err := s.verifyApplication(c, *relDoc.LinkedNamespaces.Application, &params); err != nil {
		return err
	} else {
		maps.Copy(params, appParams)
	}
	log.Info().Msgf("application verified: %s", *relDoc.LinkedNamespaces.Application)

	if _, appParams, err := s.verifyServicePrincipal(c, req.ServicePrincipalID, &params); err != nil {
		return err
	} else {
		maps.Copy(params, appParams)
	}
	log.Info().Msgf("service principal verified: %s", req.ServicePrincipalID)

	if obj, appParams, err := s.verifyGroup(c, nsID, &params); err != nil {
		return err
	} else {
		maps.Copy(params, appParams)
		doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeGroup)
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("group verified: %s", nsID)

	if ok, err := s.verifyGroupMembership(c, req.DeviceNamespaceID, nsID, &params); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%w: device is not a member of the group", common.ErrStatusBadRequest)
	}
	log.Info().Msgf("group membership verified: %s", nsID)

	return nil
}

func (s *adminServer) BeginEnrollCertificateV2(c *gin.Context, nsID uuid.UUID, templateId uuid.UUID) {

	rawReq := CertificateEnrollmentRequest{}
	if err := c.Bind(&rawReq); err != nil {
		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	req, err := rawReq.ValueByDiscriminator()
	if err != nil {
		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	switch req.(type) {
	case CertificateEnrollmentRequestDeviceLinkedServicePrincipal:
		s.processBeginEnrollCertForDASPLink(c, nsID, templateId, req.(CertificateEnrollmentRequestDeviceLinkedServicePrincipal))

	}

	respondPublicErrorMsg(c, http.StatusBadRequest, "not supported")

	/*
		if p.ValidityInMonths < 1 || p.ValidityInMonths > 120 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "validity in months must be between 1 and 120"})
			return
		}

		// check enrollment policy
		policyDoc, err := s.GetPolicyDoc(c, p.Issuer.IssuerNamespaceID, p.PolicyID)
		if err != nil {
			if common.IsAzNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"message": "no policy found for given issuer and policy id"})
				return
			}
			log.Error().Err(err).Msg("Failed to get enrollment policy")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		if policyDoc.PolicyType != PolicyTypeCertEnroll {
			c.JSON(http.StatusBadRequest, gin.H{"message": "policy is not for certificate enrollment"})
			return
		}

		// validate request against policy
		if policyDoc.CertEnroll.MaxValidityInMonths < p.ValidityInMonths {
			c.JSON(http.StatusBadRequest, gin.H{"message": "validity in months exceeds policy limit"})
			return
		}
		usageValidated := false
		for _, usage := range policyDoc.CertEnroll.AllowedUsages {
			if usage == p.Usage {
				usageValidated = true
				break
			}
		}
		if !usageValidated {
			c.JSON(http.StatusBadRequest, gin.H{"message": "usage is not allowed by policy"})
		}
	*/
}

func (s *adminServer) CompleteCertificateEnrollmentV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, certId CertIdParameter, params CompleteCertificateEnrollmentV2Params) {
	// Your code here
}
