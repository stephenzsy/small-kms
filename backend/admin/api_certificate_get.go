package admin

import "github.com/gin-gonic/gin"

func (s *adminServer) GetCertificateV2(c *gin.Context, namespaceId NamespaceIdParameter, certId CertIdParameter, params GetCertificateV2Params) {
	c.JSON(404, nil)
}
