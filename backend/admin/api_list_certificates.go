package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *adminServer) ListCertificatesV1(c *gin.Context, namespaceID uuid.UUID, params ListCertificatesV1Params) {
	if _, ok := authNamespaceRead(c, namespaceID); !ok {
		return
	}
	results := make([]CertificateRef, 0)

	c.JSON(200, results)
}
