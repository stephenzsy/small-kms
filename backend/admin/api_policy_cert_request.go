package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *adminServer) ApplyPolicyV1(c *gin.Context, namespaceID uuid.UUID, policyID uuid.UUID) {
	if _, ok := authNamespaceAdminOrSelf(c, namespaceID); !ok {
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
}
