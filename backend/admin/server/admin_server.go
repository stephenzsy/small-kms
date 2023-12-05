package adminserver

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	agentadmin "github.com/stephenzsy/small-kms/backend/admin/agent"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
)

type server struct {
	*base.BaseServer
	*profile.ProfileServer
	*agentadmin.AgentAdminServer
	*systemapp.SystemAppAdminServer
	*key.KeyAdminServer
	*cert.CertServer
	*agentadmin.AgentPushProxiedServer
}

// GetMemberGroup implements admin.ServerInterface.
func (*server) GetMemberOf(ctx echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

var _ admin.ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) (*server, error) {
	if keyAdminServer, err := key.NewServer(apiServer); err != nil {
		return nil, err
	} else if certServer, err := cert.NewServer(apiServer); err != nil {
		return nil, err
	} else {
		return &server{
			BaseServer:             base.NewBaseServer(apiServer),
			ProfileServer:          profile.NewServer(apiServer),
			AgentAdminServer:       agentadmin.NewServer(apiServer),
			SystemAppAdminServer:   systemapp.NewServer(apiServer),
			KeyAdminServer:         keyAdminServer,
			CertServer:             certServer,
			AgentPushProxiedServer: agentadmin.NewAgentPushProxiedServer(apiServer),
		}, nil
	}
}
