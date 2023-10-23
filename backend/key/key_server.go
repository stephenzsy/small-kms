package key

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
)

type server struct {
	api.APIServer
}

// GetKeyPolicy implements ServerInterface.
func (*server) GetKeyPolicy(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	panic("unimplemented")
}

// ListKeyPolicies implements ServerInterface.
func (*server) ListKeyPolicies(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	panic("unimplemented")
}

// PutKeyPolicy implements ServerInterface.
func (*server) PutKeyPolicy(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	panic("unimplemented")
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
