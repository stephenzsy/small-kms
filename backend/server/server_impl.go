package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
)

type serverImpl struct {
	admin admin.AdminServerInterface
}

func (s *serverImpl) ListCertificates(c *gin.Context, category common.CertificateCategory, params common.ListCertificatesParams) {
	if s.admin == nil {
		c.JSON(404, gin.H{"error": "Not allowed with current role"})
	}
	s.admin.ListCertificates(c, category, params)
}

func (s *serverImpl) CreateCertificate(c *gin.Context, params common.CreateCertificateParams) {
	if s.admin == nil {
		c.JSON(404, gin.H{"error": "Not allowed with current role"})
	}
	s.admin.CreateCertificate(c, params)
}

func (s *serverImpl) DownloadCertificate(c *gin.Context, id uuid.UUID, params common.DownloadCertificateParams) {
	if s.admin == nil {
		c.JSON(404, gin.H{"error": "Not allowed with current role"})
	}
	s.admin.DownloadCertificate(c, id, params)
}

func NewServerImpl() common.ServerInterface {
	serverConfig := common.NewServerConfig()
	server := serverImpl{}
	switch serverConfig.GetServerRole() {
	case common.ServerRoleAdmin:
		server.admin = admin.NewAdminServer(&serverConfig)
	}
	return &server
}
