package adminserver

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	agentadmin "github.com/stephenzsy/small-kms/backend/admin/agent"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
)

type server struct {
	*profile.ProfileServer
	*agentadmin.AgentAdminServer
	*systemapp.SystemAppAdminServer
	*key.KeyAdminServer
	*cert.CertServer
}

// GetMemberGroup implements admin.ServerInterface.
func (*server) GetMemberOf(ctx echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	panic("unimplemented")
}

var _ admin.ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		ProfileServer:        profile.NewServer(apiServer),
		AgentAdminServer:     agentadmin.NewServer(apiServer),
		SystemAppAdminServer: systemapp.NewServer(apiServer),
		KeyAdminServer:       key.NewServer(apiServer),
		CertServer:           cert.NewServer(apiServer),
	}
}
