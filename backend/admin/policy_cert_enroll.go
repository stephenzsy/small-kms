package admin

import "github.com/gin-gonic/gin"

func (s *adminServer) PutPolicyCertEnrollV1(c *gin.Context, namespaceID NamespaceID) {
	c.JSON(400, gin.H{"error": "not implemented"})
}

func (s *adminServer) GetPolicyCertEnrollV1(c *gin.Context, namespaceID NamespaceID) {
	c.JSON(400, gin.H{"error": "not implemented"})
}
