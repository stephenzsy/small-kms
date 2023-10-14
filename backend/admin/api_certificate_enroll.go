package admin

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphapplications "github.com/microsoftgraph/msgraph-sdk-go/applications"
	msgraphdevices "github.com/microsoftgraph/msgraph-sdk-go/devices"
	msgraphdirectoryobjects "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	msgraphgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	msgraphsp "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	certtemplate "github.com/stephenzsy/small-kms/backend/cert-template"
)

// extract owner name

type contextKey string

type CertificateEnrollmentClaims struct {
	jwt.RegisteredClaims
	SchemaVersion           int      `json:"ver"`
	SubjectAlternativeNames []string `json:"sans,omitempty"`
	KeyUsages               []string `json:"key_ops"`
	ExtendedUsages          []string `json:"ext_usages"`
}

const (
	graphClientContextKey contextKey = "graphClient"
)

func graphClienFromContext(c context.Context) *msgraph.GraphServiceClient {
	return c.Value(graphClientContextKey).(*msgraph.GraphServiceClient)
}

func withGraphClient(c context.Context, client *msgraph.GraphServiceClient) context.Context {
	return context.WithValue(c, graphClientContextKey, client)
}

func (s *adminServer) verifyDevice(c context.Context, objectID uuid.UUID, params *certtemplate.TemplateVarData) (msgraphmodels.Deviceable, error) {
	obj, err := graphClienFromContext(c).Devices().ByDeviceId(objectID.String()).Get(c, &msgraphdevices.DeviceItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphdevices.DeviceItemRequestBuilderGetQueryParameters{
			//Select: graph.GetProfileGraphSelectDeviceDoc(),
		},
	})
	if err != nil {
		return obj, err
	}
	params.Device = certtemplate.ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/devices/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/devices(deviceId='{%s}')", *obj.GetDeviceId()),
	}
	return obj, nil
}

func (s *adminServer) verifyApplication(c context.Context, objectID uuid.UUID, params *certtemplate.TemplateVarData) (msgraphmodels.Applicationable, error) {
	obj, err := graphClienFromContext(c).Applications().ByApplicationId(objectID.String()).Get(c, &msgraphapplications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphapplications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, err
	}
	params.Application = certtemplate.ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/applications/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/applications(appId='{%s}')", *obj.GetAppId()),
	}
	return obj, nil
}

func (s *adminServer) verifyServicePrincipal(c context.Context, objectID uuid.UUID, params *certtemplate.TemplateVarData) (msgraphmodels.ServicePrincipalable, error) {
	obj, err := graphClienFromContext(c).ServicePrincipals().ByServicePrincipalId(objectID.String()).Get(c, &msgraphsp.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphsp.ServicePrincipalItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, err
	}
	params.ServicePrincipal = certtemplate.ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals(appId='{%s}')", *obj.GetAppId()),
	}
	return obj, nil
}

func (s *adminServer) verifyGroup(c context.Context, objectID uuid.UUID, params *certtemplate.TemplateVarData) (msgraphmodels.Groupable, error) {
	obj, err := graphClienFromContext(c).Groups().ByGroupId(objectID.String()).Get(c, &msgraphgroups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphgroups.GroupItemRequestBuilderGetQueryParameters{
			//Select: graph.GetProfileGraphSelectGroupDoc(),
		},
	})
	if err != nil {
		return obj, err
	}
	params.Group = certtemplate.ResourceTemplateVarData{
		ID:  *obj.GetId(),
		URI: fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s", *obj.GetId()),
	}
	return obj, nil
}

func (s *adminServer) verifyGroupMembership(c context.Context, objectID uuid.UUID, groupID uuid.UUID) (bool, error) {
	requestBody := msgraphdirectoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	requestBody.SetGroupIds([]string{groupID.String()})

	checkMemberGroups, err := graphClienFromContext(c).DirectoryObjects().ByDirectoryObjectId(objectID.String()).CheckMemberGroups().Post(c, requestBody, nil)
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

/*

func (s *adminServer) processBeginEnrollCertForDASPLink(c context.Context, nsID uuid.UUID, templateID uuid.UUID, req CertificateEnrollmentRequestDeviceLinkedServicePrincipal) (*PendingCertDoc, error) {
	log.Info().Msgf("enroll cert for dasp link - begin: %s", req.DeviceLinkID)
	defer log.Info().Msgf("enroll cert for dasp link - end: %s", req.DeviceLinkID)

	// first check if AppID match
	authCtx, ok := auth.GetAuthIdentity(c)
	if !ok {
		return nil, fmt.Errorf("%w: no auth identity", common.ErrStatusUnauthorized)
	}
	requesterID := authCtx.ClientPrincipalID()
	if requesterID == uuid.Nil {
		return nil, fmt.Errorf("%w: client principal id is required", common.ErrStatusUnauthorized)
	}
	if req.AppID == uuid.Nil {
		return nil, fmt.Errorf("%w: application id is nil", common.ErrStatusBadRequest)
	}
	claimedAppId := authCtx.AppIDClaim()
	if claimedAppId != req.AppID {
		return nil, fmt.Errorf("%w: appid is required in the authorization claim", common.ErrStatusUnauthorized)
	}

	// get template
	// template, err := certtemplate.LoadCertifictateTemplate(c, nsID, templateID)
	// if err != nil {
	// 	return bad(err)
	// }
	// // verify template is active
	// if !template.IsEnabled() {
	// 	return bad(fmt.Errorf("%w: template is not enabled", common.ErrStatusBadRequest))
	// }

	// graph client in context

	// look up relDoc
	if req.DeviceLinkID != common.GetCanonicalNamespaceRelationID(req.DeviceNamespaceID, common.NSRelNameDASPLink) {
		return nil, fmt.Errorf("%w: device link id is invalid", common.ErrStatusBadRequest)
	}
	// relDoc, err := s.readNsRel(c, req.DeviceNamespaceID, req.DeviceLinkID)
	// if err != nil {
	// 	return nil, err
	// }
	// if relDoc.Status != NsRelStatusEnabled {
	// 	return nil, fmt.Errorf("%w: device link is not enabled", common.ErrStatusBadRequest)
	// }
	// if deviceID, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.Device); !nonNil || deviceID != req.DeviceNamespaceID {
	// 	return nil, fmt.Errorf("%w: device object id invalid", common.ErrStatusBadRequest)
	// }
	// if spID, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.ServicePrincipal); !nonNil || spID != req.ServicePrincipalID {
	// 	return nil, fmt.Errorf("%w: service principal id invalid", common.ErrStatusBadRequest)
	// }
	// if appID, nonNil := utils.NonNilUUID(relDoc.Attributes.AppID); !nonNil || appID != req.AppID {
	// 	return nil, fmt.Errorf("%w: application client id does invalid", common.ErrStatusBadRequest)
	// }
	// if _, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.Application); !nonNil {
	// 	return nil, fmt.Errorf("%w: application object id is nil", common.ErrStatusBadRequest)
	// }
	log.Info().Msgf("link doc loaded and verified: %s", req.DeviceLinkID)

	// prep parameters
	params := certtemplate.TemplateVarData{}
	// verify against ms graph
	if _, err := s.verifyDevice(c, req.DeviceNamespaceID, &params); err != nil {
		return nil, err
	} else {
		// doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeDevice)
		// if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
		// 	return nil, err
		// }
	}
	log.Info().Msgf("device verified: %s", req.DeviceNamespaceID)

	// if _, err := s.verifyApplication(c, *relDoc.LinkedNamespaces.Application, &params); err != nil {
	// 	return nil, err
	// }
	// log.Info().Msgf("application verified: %s", *relDoc.LinkedNamespaces.Application)

	// if _, err := s.verifyServicePrincipal(c, req.ServicePrincipalID, &params); err != nil {
	// 	return nil, err
	// }
	log.Info().Msgf("service principal verified: %s", req.ServicePrincipalID)

	if _, err := s.verifyGroup(c, nsID, &params); err != nil {
		return nil, err
		// } else {
		// 	doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeGroup)
		// 	if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
		// 		return nil, err
		// 	}
	}
	log.Info().Msgf("group verified: %s", nsID)

	if ok, err := s.verifyGroupMembership(c, req.DeviceNamespaceID, nsID); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("%w: device is not a member of the group", common.ErrStatusBadRequest)
	}
	log.Info().Msgf("group membership verified: %s", nsID)

	// cert, err := template.CreateCertWithVariables(params)
	// if err != nil {
	// 	return bad(err)
	// }
	// claimsEncoded, err := encodeJwtJsonSegment(cert.TODO()) // get jwt claims
	// if err != nil {
	// 	return nil, err
	// }

	// store doc
	// pCertDoc := newPendingCertDoc(certID, cert, claimsEncoded, templateDoc, req.ServicePrincipalID, requesterID)
	// if err = kmsdoc.AzCosmosCreate(c, s.AzCosmosContainerClient(), &pCertDoc); err != nil {
	// 	return nil, err
	// }
	// pCertDoc := cert.TODO().(*PendingCertDoc)
	// log.Info().Msgf("pending document stored: %s", nsID)

	// return pCertDoc, nil
	return nil, nil
}
*/
// func (s *adminServer) BeginEnrollCertificateV2(c *gin.Context, nsID uuid.UUID, templateId uuid.UUID) {

// 	rawReq := CertificateEnrollmentRequest{}
// 	if err := c.Bind(&rawReq); err != nil {
// 		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	_, err := rawReq.ValueByDiscriminator()
// 	if err != nil {
// 		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	//serviceContext := common.CreateServiceContext(c, s.AzCosmosContainerClient())

// 	var pCertDoc *PendingCertDoc
// 	//var responseNsType NamespaceTypeShortName
// 	//	switch req := req.(type) {
// 	// case CertificateEnrollmentRequestDeviceLinkedServicePrincipal:
// 	// 	//pCertDoc, err = s.processBeginEnrollCertForDASPLink(serviceContext, nsID, templateId, req)
// 	// 	responseNsType = NSTypeServicePrincipal
// 	// }
// 	if err != nil {
// 		common.RespondError(c, err)
// 		return
// 	}

// 	if pCertDoc == nil {
// 		respondPublicErrorMsg(c, http.StatusBadRequest, "not supported")
// 	}
// 	//c.JSON(http.StatusCreated, pCertDoc.toReceipt(responseNsType))
// }

// func (s *adminServer) CompleteCertificateEnrollmentV2(c *gin.Context, nsID uuid.UUID, certID uuid.UUID, params CompleteCertificateEnrollmentV2Params) {
// 	if _, ok := authNamespaceAdminOrSelf(c, nsID); !ok {
// 		return
// 	}

// 	req := new(CertificateEnrollmentReplyFinalize)
// 	if err := c.Bind(req); err != nil {
// 		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	pCertDoc, err := s.readPendingCertDoc(c, nsID, kmsdoc.NewKmsDocID(kmsdoc.DocTypePendingCert, certID))
// 	if err != nil {
// 		common.RespondError(c, err)
// 	}

// 	parser := jwt.NewParser()

// 	completeToken := req.JwtHeader + "." + pCertDoc.JWT[1] + "." + req.JwtSignature

// 	pubKey := rsa.PublicKey{}
// 	/*err = req.PublicKey.populateRsaPublicKey(&pubKey)
// 	if err != nil {
// 		log.Warn().Err(err).Msgf("failed to parse public key: %s", *req.PublicKey.KeyID)
// 		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
// 		return
// 	}*/
// 	_, err = parser.Parse(completeToken, func(_ *jwt.Token) (interface{}, error) {
// 		return &pubKey, nil
// 	})
// 	if err != nil {
// 		log.Warn().Err(err).Msgf("failed to parse jwt token: %s", completeToken)
// 		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	decodedClaims, err := base64.RawStdEncoding.DecodeString(pCertDoc.JWT[1])
// 	if err != nil {
// 		respondInternalError(c, err, "failed to decode jwt claims")
// 		return
// 	}
// 	claims := CertificateEnrollmentClaims{}
// 	if err := json.Unmarshal(decodedClaims, &claims); err != nil {
// 		respondInternalError(c, err, "failed to parse jwt claims")
// 		return
// 	}
// 	cert := new(x509.Certificate)
// 	pCertDoc.populateCertificate(cert)

// 	//s.createCertificateFromTemplateWithCert(c, nsID, pCertDoc.TemplateDoc, cert, certID)

// 	c.JSON(http.StatusOK, nil)
// }

/*
func (p *models.JwkProperties) populateRsaPublicKey(k *rsa.PublicKey) error {
	if p == nil {
		return nil
	}
	decE, err := base64.RawURLEncoding.DecodeString(*p.E)
	if err != nil {
		return err
	}
	decN, err := base64.RawURLEncoding.DecodeString(*p.N)
	if err != nil {
		return err
	}

	k.E = int(new(big.Int).SetBytes(decE).Int64())
	k.N = new(big.Int).SetBytes(decN)
	return nil
}
*/
