package key

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
)

type server struct {
	api.APIServer
}

// ListKeySpecs implements ServerInterface.
func (*server) ListKeySpecs(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	panic("unimplemented")
}

// GetKeySpec implements ServerInterface.
func (*server) GetKeySpec(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	panic("unimplemented")
}

// PutKeySpec implements ServerInterface.
func (*server) PutKeySpec(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	panic("unimplemented")
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
