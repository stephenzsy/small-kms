package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
)

func (s *adminServer) MyHasPermissionV1(c *gin.Context, permissionKey NamespacePermissionKey) {
	s.HasPermissionV1(c, auth.CallerPrincipalId(c), permissionKey)
}

func (s *adminServer) HasPermissionV1(c *gin.Context, namespaceId uuid.UUID, permissionKey NamespacePermissionKey) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceId); !ok {
		return
	}

	switch permissionKey {
	case AllowEnrollDeviceCertificate:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid permission key %s", string(permissionKey))})
		return
	}
	docs, err := s.queryNsRelHasPermission(c, namespaceId, permissionKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get permissions")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
	}
	r := make([]NamespaceRef, len(docs))
	for i, doc := range docs {
		r[i].ID = doc.ID.GetUUID()
		r[i].DisplayName = doc.DisplayName
	}

	c.JSON(http.StatusOK, r)
}
