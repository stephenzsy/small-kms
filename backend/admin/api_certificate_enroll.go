package admin

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
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

func (s *adminServer) verifyDevice(c context.Context, objectID uuid.UUID, params *TemplateVarData) (msgraphmodels.Deviceable, error) {
	obj, err := graphClienFromContext(c).Devices().ByDeviceId(objectID.String()).Get(c, &msgraphdevices.DeviceItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphdevices.DeviceItemRequestBuilderGetQueryParameters{
			Select: graph.GetProfileGraphSelectDeviceDoc(),
		},
	})
	if err != nil {
		return obj, err
	}
	params.Device = ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/devices/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/devices(deviceId='{%s}')", *obj.GetDeviceId()),
	}
	return obj, nil
}

func (s *adminServer) verifyApplication(c context.Context, objectID uuid.UUID, params *TemplateVarData) (msgraphmodels.Applicationable, error) {
	obj, err := graphClienFromContext(c).Applications().ByApplicationId(objectID.String()).Get(c, &msgraphapplications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphapplications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, err
	}
	params.Application = ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/applications/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/applications(appId='{%s}')", *obj.GetAppId()),
	}
	return obj, nil
}

func (s *adminServer) verifyServicePrincipal(c context.Context, objectID uuid.UUID, params *TemplateVarData) (msgraphmodels.ServicePrincipalable, error) {
	obj, err := graphClienFromContext(c).ServicePrincipals().ByServicePrincipalId(objectID.String()).Get(c, &msgraphsp.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphsp.ServicePrincipalItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId"},
		},
	})
	if err != nil {
		return obj, err
	}
	params.ServicePrincipal = ResourceTemplateVarData{
		ID:     *obj.GetId(),
		URI:    fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals/%s", *obj.GetId()),
		AltURI: fmt.Sprintf("https://graph.microsoft.com/v1.0/servicePrincipals(appId='{%s}')", *obj.GetAppId()),
	}
	return obj, nil
}

func (s *adminServer) verifyGroup(c context.Context, objectID uuid.UUID, params *TemplateVarData) (msgraphmodels.Groupable, error) {
	obj, err := graphClienFromContext(c).Groups().ByGroupId(objectID.String()).Get(c, &msgraphgroups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &msgraphgroups.GroupItemRequestBuilderGetQueryParameters{
			Select: graph.GetProfileGraphSelectGroupDoc(),
		},
	})
	if err != nil {
		return obj, err
	}
	params.Group = ResourceTemplateVarData{
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

func processTemplate(tmplStr string, data *TemplateVarData) string {
	tmplStr = strings.TrimSpace(tmplStr)
	tmpl, err := parseCertificateRequestTemplate(tmplStr)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to parse template: %s", tmplStr)
		return ""
	}
	if tmpl == nil {
		return tmplStr
	}
	transformed, err := executeTemplate(tmpl, data)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to execute template: %s", tmplStr)
		return ""
	}
	return transformed
}

func processCertificateEnrollmentClaims(template *CertificateTemplateDoc, data *TemplateVarData) (certID uuid.UUID, claims CertificateEnrollmentClaims, err error) {
	claims.SchemaVersion = 1
	var cert *x509.Certificate
	cert, certID, err = prepareUnsignedCertificateFromTemplate(NamespaceTypeShortName(""), uuid.Nil, template, data)
	if err != nil {
		return
	}
	claims.Subject = cert.Subject.String()
	claims.NotBefore = jwt.NewNumericDate(cert.NotBefore)
	claims.ExpiresAt = jwt.NewNumericDate(cert.NotAfter)
	claims.ID = certID.String()
	if len(cert.EmailAddresses) > 0 {
		claims.SubjectAlternativeNames = append(claims.SubjectAlternativeNames, cert.EmailAddresses...)
	}
	if len(cert.URIs) > 0 {
		claims.SubjectAlternativeNames = append(claims.SubjectAlternativeNames, utils.MapSlices(cert.URIs,
			func(uri *url.URL) string { return uri.String() })...)
	}
	if cert.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
		claims.KeyUsages = append(claims.KeyUsages, "sign", "verify")
	}
	if cert.KeyUsage&x509.KeyUsageDataEncipherment != 0 {
		claims.KeyUsages = append(claims.KeyUsages, "encrypt", "decrypt")
	}
	if cert.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
		claims.KeyUsages = append(claims.KeyUsages, "wrapKey", "unwrapKey")
	}
	claims.ExtendedUsages = utils.FilterSlice(utils.MapSlices(cert.ExtKeyUsage, func(u x509.ExtKeyUsage) string {
		switch u {
		case x509.ExtKeyUsageServerAuth:
			return "1.3.6.1.5.5.7.3.1"
		case x509.ExtKeyUsageClientAuth:
			return "1.3.6.1.5.5.7.3.2"
		}
		return ""
	}), func(u string) bool { return u != "" })
	return
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

	// get template
	templateDoc, err := s.readCertificateTemplateDoc(c, nsID, templateId)
	if err != nil {
		return err
	}

	// graph client in contextF
	if graphClient, err := s.msGraphClient(c); err != nil {
		return err
	} else {
		c = withGraphClient(c, graphClient)
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
		return fmt.Errorf("%w: service principal id invalid", common.ErrStatusBadRequest)
	}
	if appID, nonNil := utils.NonNilUUID(relDoc.Attributes.AppID); !nonNil || appID != req.AppID {
		return fmt.Errorf("%w: application client id does invalid", common.ErrStatusBadRequest)
	}
	if _, nonNil := utils.NonNilUUID(relDoc.LinkedNamespaces.Application); !nonNil {
		return fmt.Errorf("%w: application object id is nil", common.ErrStatusBadRequest)
	}
	log.Info().Msgf("link doc loaded and verified: %s", req.DeviceLinkID)

	// prep parameters
	params := TemplateVarData{}
	// verify against ms graph
	if obj, err := s.verifyDevice(c, req.DeviceNamespaceID, &params); err != nil {
		return err
	} else {
		doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeDevice)
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("device verified: %s", req.DeviceNamespaceID)

	if _, err := s.verifyApplication(c, *relDoc.LinkedNamespaces.Application, &params); err != nil {
		return err
	}
	log.Info().Msgf("application verified: %s", *relDoc.LinkedNamespaces.Application)

	if _, err := s.verifyServicePrincipal(c, req.ServicePrincipalID, &params); err != nil {
		return err
	}
	log.Info().Msgf("service principal verified: %s", req.ServicePrincipalID)

	if obj, err := s.verifyGroup(c, nsID, &params); err != nil {
		return err
	} else {
		doc := s.graphService.NewGraphProfileDocWithType(s.TenantID(), obj, graph.MsGraphOdataTypeGroup)
		if err := kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("group verified: %s", nsID)

	if ok, err := s.verifyGroupMembership(c, req.DeviceNamespaceID, nsID); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%w: device is not a member of the group", common.ErrStatusBadRequest)
	}
	log.Info().Msgf("group membership verified: %s", nsID)

	certID, claims, err := processCertificateEnrollmentClaims(templateDoc, &params)
	if err != nil {
		return err
	}
	marshalled, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	base64.RawURLEncoding.EncodeToString(marshalled)

	log.Info().Msgf("certificate enrollment claims processed: %s, %s", certID, base64.RawURLEncoding.EncodeToString(marshalled))

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

	switch req := req.(type) {
	case CertificateEnrollmentRequestDeviceLinkedServicePrincipal:
		err = s.processBeginEnrollCertForDASPLink(c, nsID, templateId, req)

	}
	if err != nil {
		common.RespondError(c, err)
		return
	}

	respondPublicErrorMsg(c, http.StatusBadRequest, "not supported")

}

func (s *adminServer) CompleteCertificateEnrollmentV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, certId CertIdParameter, params CompleteCertificateEnrollmentV2Params) {
	// Your code here
}
