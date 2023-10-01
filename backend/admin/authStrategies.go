package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
)

func authNamespaceAdminOrSelf(c *gin.Context, namespaceID uuid.UUID) (principalID uuid.UUID, ok bool) {
	if identity, ok := auth.GetAuthIdentity(c); ok {
		principalID = identity.ClientPrincipalID()
		if identity.HasAdminRole() {
			return principalID, true
		} else if ok && (principalID == namespaceID) {
			return principalID, true
		}
	}

	respondPublicErrorMsg(c, http.StatusForbidden, fmt.Sprintf("caller %s does not access to namespace: %s", principalID, namespaceID))
	ok = false
	return
}

func authAdminOnly(c *gin.Context) bool {
	if identity, ok := auth.GetAuthIdentity(c); ok {
		if identity.HasAdminRole() {
			return true
		}
	}
	respondPublicErrorMsg(c, http.StatusForbidden, "caller does not have admin addess")
	return false
}

func authGroupOwnerWithLinkedDevice(c *gin.Context, nsID uuid.UUID, onSyncDoc func(uuid.UUID) (*DirectoryObjectDoc, error), onBehalfOf uuid.UUID) bool {
	identity, ok := auth.GetAuthIdentity(c)
	if !ok {
		respondPublicErrorMsg(c, http.StatusUnauthorized, "caller is not authenticated")
		return false
	}

	groupDoc, err := onSyncDoc(nsID)
	if err != nil {
		if common.IsAzNotFound(err) || common.IsGraphODataErrorNotFound(err) {
			respondPublicErrorMsg(c, http.StatusNotFound, "group not found")
			return false
		}
		respondInternalError(c, err, "failed to sync group")
	}
	if groupDoc.OdataType != "#microsoft.graph.group" {
		respondPublicErrorMsg(c, http.StatusBadRequest, "target is not a group")
		return false
	}

	stringGraphCheckIDs := make([]string, 0, 4)
	if onBehalfOf != uuid.Nil {
	}

	onSyncDoc(nsID)
}
