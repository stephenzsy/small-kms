package profile

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type server struct {
	api.APIServer
}

const namespaceIDRootCA = "root-ca"

var namespaceIdentifierRootCA = base.StringIdentifier(namespaceIDRootCA)

// ListRootCAs implements ServerInterface.
func (s *server) ListRootCAs(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	result, err := listProfiles(c, namespaceIdentifierRootCA, base.ProfileResourceKindRootCA)
	if err != nil {
		return err
	}
	if result == nil {
		result = []*ProfileRef{}
	}
	return c.JSON(200, result)
}

// GetRootCA implements ServerInterface.
func (*server) GetRootCA(ctx echo.Context, namespaceIdentifier base.Identifier) error {
	panic("unimplemented")
}

// PutRootCA implements ServerInterface.
func (*server) PutRootCA(ctx echo.Context, namespaceIdentifier base.Identifier) error {
	panic("unimplemented")
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
