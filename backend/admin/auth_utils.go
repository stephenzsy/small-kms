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
		respondPublicError(c, http.StatusForbidden, fmt.Errorf("caller %s does not access to namespace: %s", callerID, namespaceID))
		ok = false
	}
	return
}

func authNamespaceRead(c *gin.Context, namespaceID uuid.UUID) (callerID uuid.UUID, ok bool) {
	callerID = auth.CallerPrincipalId(c)
	ok = true
	if !IsCANamespace(namespaceID) && callerID != namespaceID {
		respondPublicError(c, http.StatusForbidden, fmt.Errorf("caller %s does not have access to namespace %s", callerID, namespaceID))
		ok = false
		return
	}
	return
}
