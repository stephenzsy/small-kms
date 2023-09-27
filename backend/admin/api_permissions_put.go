package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

func (s *adminServer) PutPermissionsV1(c *gin.Context, namespaceId uuid.UUID, objectId uuid.UUID) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	p := new(NamespacePermissions)
	if err := c.BindJSON(p); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "invalid request body"})
	}

	// verify namespaces of user to device
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
	if profile.ObjectType != NamespaceTypeMsGraphUser {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("namespace %s is not a user", namespaceId)})
		return
	}

	objectProfile, status, err := s.RegisterNamespaceProfile(c, objectId)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.Error().Err(err).Msg("Failed to register graph object")
			c.JSON(status, gin.H{"message": "internal error"})
		} else {
			c.JSON(status, gin.H{"message": err.Error()})
		}
		return
	}
	if objectProfile.ObjectType != NamespaceTypeMsGraphDevice && !s.skipDeviceCheck {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("target namespace %s is not a device", namespaceId)})
		return
	}

	doc := new(NsRelDoc)
	doc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, objectId)
	doc.NamespaceID = namespaceId
	doc.AllowEnrollDeviceCertificate = p.AllowEnrollDeviceCertificate
	doc.DisplayName = objectProfile.DisplayName

	if err := kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, doc); err != nil {
		log.Error().Err(err).Msg("Failed to store permission")
		c.JSON(status, gin.H{"message": "internal error"})
		return
	}

	c.JSON(http.StatusOK, p)
}
