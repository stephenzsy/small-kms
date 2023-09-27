package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
)

func authNamespaceAdminOrSelf(c *gin.Context, namespaceID uuid.UUID) (callerID uuid.UUID, ok bool) {
	callerID = auth.CallerPrincipalId(c)
	ok = true
	if !auth.CallerPrincipalHasAdminRole(c) && callerID != namespaceID {
		respondPublicErrorMsg(c, http.StatusForbidden, fmt.Sprintf("caller %s does not access to namespace: %s", callerID, namespaceID))
		ok = false
	}
	return
}

func authAdminOnly(c *gin.Context) bool {
	if !auth.CallerPrincipalHasAdminRole(c) {
		respondPublicErrorMsg(c, http.StatusForbidden, "caller does not have admin addess")
		return false
	}
	return true
}

func authNamespaceRead(c *gin.Context, namespaceID uuid.UUID) (callerID uuid.UUID, ok bool) {
	callerID = auth.CallerPrincipalId(c)
	ok = true
	if !IsCANamespace(namespaceID) && !auth.CallerPrincipalHasAdminRole(c) && callerID != namespaceID {
		respondPublicErrorMsg(c, http.StatusForbidden, fmt.Sprintf("caller %s does not have access to namespace %s", callerID, namespaceID))
		ok = false
		return
	}
	return
}
