package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/stephenzsy/small-kms/backend/common"
)

type AdminServerInterface interface {
	ListCACertificates(c *gin.Context, params common.ListCACertificatesParams)
	CreateCertificate(c *gin.Context, params common.CreateCertificateParams)
}

type adminServer struct {
	config common.ServerConfig
}

func NewAdminServer(c common.ServerConfig) AdminServerInterface {
	return &adminServer{config: c}
}
