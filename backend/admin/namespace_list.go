package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
)

func (s *adminServer) ListNamespacesV1(c *gin.Context, namespaceType NamespaceType) {
	if !auth.CallerPrincipalHasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admin can list name spaces"})
		return
	}

	switch namespaceType {
	case NamespaceTypeMsGraphServicePrincipal:
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
		results[i] = NamespaceRef{
			ID:                   item.ID.GetUUID(),
			DisplayName:          item.DisplayName,
			ObjectType:           NamespaceType(item.OdataType),
			UserPrincipalName:    item.UserPrincipalName,
			ServicePrincipalType: item.ServicePrincipalType,
		}
	}
	c.JSON(http.StatusOK, results)
}
