package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/common"
)

type AdminServerInterface interface {
	ListCertificates(c *gin.Context, category common.CertificateCategory, params common.ListCertificatesParams)
	CreateCertificate(c *gin.Context, params common.CreateCertificateParams)
	DownloadCertificate(c *gin.Context, id uuid.UUID, params common.DownloadCertificateParams)
}

type adminServer struct {
	config common.ServerConfig
}

func NewAdminServer(c common.ServerConfig) AdminServerInterface {
	return &adminServer{config: c}
}
