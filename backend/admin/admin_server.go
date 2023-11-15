package admin

import (
	agentadmin "github.com/stephenzsy/small-kms/backend/admin/agent"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/api"
)

type server struct {
	*agentadmin.AgentAdminServer
	*systemapp.SystemAppAdminServer
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		AgentAdminServer: agentadmin.NewServer(apiServer),
	}
}
