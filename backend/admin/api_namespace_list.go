package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
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
			NamespaceID:          directoryID,
			ID:                   item.ID.GetUUID(),
			DisplayName:          item.DisplayName,
			ObjectType:           NamespaceType(item.OdataType),
			UserPrincipalName:    item.UserPrincipalName,
			ServicePrincipalType: item.ServicePrincipalType,
			Updated:              item.Updated,
			UpdatedBy:            item.UpdatedBy,
		}
	}
	c.JSON(http.StatusOK, results)
}
