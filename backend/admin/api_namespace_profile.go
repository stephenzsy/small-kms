package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func getBuiltInCaIntNamespaceRefs() []NamespaceRef {
	return []NamespaceRef{
		{NamespaceID: wellKnownNamespaceID_IntCaIntranet, ID: wellKnownNamespaceID_IntCaIntranet, DisplayName: "Intermediate CA - Intranet", ObjectType: NamespaceTypeBuiltInCaInt},
		{NamespaceID: testNamespaceID_IntCA, ID: testNamespaceID_IntCA, DisplayName: "Test Intermediate CA", ObjectType: NamespaceTypeBuiltInCaInt},
	}
}

func (s *adminServer) ListNamespacesV1(c *gin.Context, namespaceType NamespaceType) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin can list name spaces"})
		return
	}

	switch namespaceType {
	case NamespaceTypeBuiltInCaInt:
		c.JSON(http.StatusOK, getBuiltInCaIntNamespaceRefs())
		return
	case NamespaceTypeMsGraphGroup,
		NamespaceTypeMsGraphServicePrincipal:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "namespace type not supported"})
		return
	}

	list, err := s.ListDirectoryObjectByType(c, namespaceType)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get list of directory objects")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	results := make([]NamespaceRef, len(list))
	for i, item := range list {
		item.PopulateNamespaceRef(&results[i])
	}
	c.JSON(http.StatusOK, results)
}

func (s *adminServer) RegisterNamespaceProfile(c *gin.Context, objectID uuid.UUID) (*NamespaceProfile, int, error) {
	dirObj, err := s.msGraphClient.DirectoryObjects().ByDirectoryObjectId(objectID.String()).Get(c, nil)
	if err != nil {
		if common.IsGraphODataErrorNotFound(err) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	doc := new(DirectoryObjectDoc)
	doc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID)
	doc.NamespaceID = directoryID
	doc.OdataType = *dirObj.GetOdataType()
	switch doc.OdataType {
	case "#microsoft.graph.user":
		if userObj, ok := dirObj.(msgraphmodels.Userable); ok {
			doc.DisplayName = *userObj.GetDisplayName()
			doc.UserPrincipalName = userObj.GetUserPrincipalName()
		}
	case "#microsoft.graph.servicePrincipal":
		if spObj, ok := dirObj.(msgraphmodels.ServicePrincipalable); ok {
			doc.DisplayName = *spObj.GetDisplayName()
			doc.ServicePrincipalType = spObj.GetServicePrincipalType()
		}
	case "#microsoft.graph.group":
		if gObj, ok := dirObj.(msgraphmodels.Groupable); ok {
			doc.DisplayName = *gObj.GetDisplayName()
		}
	case "#microsoft.graph.device":
		if dObj, ok := dirObj.(msgraphmodels.Deviceable); ok {
			doc.DisplayName = *dObj.GetDisplayName()
			doc.DeviceID = dObj.GetDeviceId()
			doc.OperatingSystem = dObj.GetOperatingSystem()
			doc.OperatingSystemVersion = dObj.GetOperatingSystemVersion()
			doc.DeviceOwnership = dObj.GetDeviceOwnership()
			doc.IsCompliant = dObj.GetIsCompliant()
		}
	default:
		return nil, http.StatusBadRequest, fmt.Errorf("graph object type (%s) not supported", doc.OdataType)
	}

	err = kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, doc)
	if err != nil {
		return nil, http.StatusInternalServerError, err

	}

	nsProfile := new(NamespaceProfile)
	doc.PopulateNamespaceProfile(nsProfile)

	return nsProfile, http.StatusOK, nil
}

func (s *adminServer) RegisterNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin can register namespaces"})
		return
	}
	if namespaceId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No namespace id specified"})
		return
	}
	profile, status, err := s.RegisterNamespaceProfile(c, namespaceId)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.Error().Err(err).Msg("Failed to register graph object")
			c.JSON(status, gin.H{"message": "internal error"})
		} else {
			c.JSON(status, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (s *adminServer) GetNamespaceProfile(c context.Context, namespaceId uuid.UUID) (*NamespaceProfile, error) {
	doc, err := s.GetDirectoryObjectDoc(c, namespaceId)
	if common.IsAzNotFound(err) {
		return nil, nil
	}
	nsProfile := new(NamespaceProfile)
	doc.PopulateNamespaceProfile(nsProfile)
	return nsProfile, nil
}

func (s *adminServer) GetNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
		return
	}

	profile, err := s.GetNamespaceProfile(c, namespaceId)
	if err != nil {
		log.Error().Err(err).Msg("Internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, profile)
}
