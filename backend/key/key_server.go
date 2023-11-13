package key

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
)

type server struct {
	api.APIServer
}

// GenerateKey implements ServerInterface.
func (*server) GenerateKey(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, resourceId base.ID) error {
	panic("unimplemented")
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
