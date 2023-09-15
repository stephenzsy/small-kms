package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
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

func (s *adminServer) RegisterNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
		return
	}
	if namespaceId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No namespace id specified"})
		return
	}
	dirObj, err := s.msGraphClient.DirectoryObjects().ByDirectoryObjectId(namespaceId.String()).Get(c, nil)
	if err != nil {
		if odErr, ok := err.(*odataerrors.ODataError); ok {
			if odErr.ResponseStatusCode == http.StatusNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "graph object not found"})
				return
			}
		}
		log.Error().Err(err).Msg("Failed to get graph object")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	doc := new(DirectoryObjectDoc)
	doc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, namespaceId)
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
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("graph object type (%s) not supported", doc.OdataType)})
		return
	}

	err = kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	nsProfile := NamespaceProfile{}
	doc.PopulateNamespaceRef(&nsProfile)

	c.JSON(http.StatusOK, &nsProfile)
}

func (s *adminServer) GetNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
		return
	}

	doc, err := s.GetDirectoryObjectDoc(c, namespaceId)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		log.Error().Err(err).Msg("Internal error")
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	nsProfile := NamespaceProfile{}
	doc.PopulateNamespaceRef(&nsProfile)

	c.JSON(http.StatusOK, &nsProfile)
}
