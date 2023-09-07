package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) RegisterNamespaceProfileV1(c *gin.Context, namespaceId uuid.UUID) {
	callerId, ok := authNamespaceAdminOrSelf(c, namespaceId)
	if !ok && callerId != namespaceId {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin or self can sync graph"})
		return
	}
	if namespaceId == uuid.Nil {
		namespaceId = callerId
	}
	if namespaceId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No namespace id specified"})
		return
	}
	dirObj, err := s.msGraphClient.DirectoryObjects().ByDirectoryObjectIdString(namespaceId.String()).Get(c, nil)
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
	callerId, ok := authNamespaceAdminOrSelf(c, namespaceId)
	if !ok && callerId != namespaceId {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin or self can sync graph"})
		return
	}

	doc, err := s.GetDirectoryObjectDoc(c, namespaceId)
	if err != nil {
		if common.IsAzNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
	}

	nsProfile := NamespaceProfile{}
	doc.PopulateNamespaceRef(&nsProfile)

	c.JSON(http.StatusOK, &nsProfile)
}
