package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/stephenzsy/small-kms/backend/common"
)

type AdminServerInterface interface {
	ListCACertificates(c *gin.Context, params common.ListCACertificatesParams)
	CreateCACertificate(c *gin.Context, id string, params common.CreateCACertificateParams)
}

type adminServer struct {
	config common.ServerConfig
}

func NewAdminServer(c common.ServerConfig) AdminServerInterface {
	return &adminServer{config: c}
}
