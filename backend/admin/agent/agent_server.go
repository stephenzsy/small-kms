package agentadmin

import (
	"github.com/stephenzsy/small-kms/backend/api"
)

type AgentAdminServer struct {
	api.APIServer
}

func NewServer(apiServer api.APIServer) *AgentAdminServer {
	return &AgentAdminServer{
		APIServer: apiServer,
	}
}
