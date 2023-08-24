package server

import (
	"github.com/gin-gonic/gin"

	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
)

type serverImpl struct {
	admin admin.AdminServerInterface
}

func (s *serverImpl) ListCACertificates(c *gin.Context, params common.ListCACertificatesParams) {
	if s.admin == nil {
		c.JSON(404, gin.H{"error": "Not allowed"})
	}
	s.admin.ListCACertificates(c, params)
}

func (s *serverImpl) CreateCACertificate(c *gin.Context, id string, params common.CreateCACertificateParams) {
	if s.admin == nil {
		c.JSON(404, gin.H{"error": "Not allowed"})
	}
	s.admin.CreateCACertificate(c, id, params)
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
