package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ListCACertificates(c *gin.Context, params common.ListCACertificatesParams) {
	c.JSON(200, &common.CertificateRefs{})
}

func (s *adminServer) CreateCACertificate(c *gin.Context, id string, params common.CreateCACertificateParams) {

}
