package admin

import (
	echo "github.com/labstack/echo/v4"
	agentadmin "github.com/stephenzsy/small-kms/backend/admin/agent"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
)

type server struct {
	*agentadmin.AgentAdminServer
	*systemapp.SystemAppAdminServer
	*key.KeyAdminServer
	*cert.CertServer
}

// GetAgentConfig implements ServerInterface.
func (*server) GetAgentConfig(ctx echo.Context, namespaceId string) error {
	panic("unimplemented")
}

// PutAgentConfig implements ServerInterface.
func (*server) PutAgentConfig(ctx echo.Context, namespaceId string) error {
	panic("unimplemented")
}

// GetCertificatePolicy implements ServerInterface.
func (*server) GetCertificatePolicy(ctx echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	panic("unimplemented")
}

// PutCertificatePolicy implements ServerInterface.
func (*server) PutCertificatePolicy(ctx echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	panic("unimplemented")
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		AgentAdminServer:     agentadmin.NewServer(apiServer),
		SystemAppAdminServer: systemapp.NewServer(apiServer),
		KeyAdminServer:       key.NewServer(apiServer),
		CertServer:           cert.NewServer(apiServer),
	}
}
