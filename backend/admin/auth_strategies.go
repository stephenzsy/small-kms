package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
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
