package managedapp

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type server struct {
	api.APIServer
}

// ListManagedApps implements ServerInterface.
func (s *server) ListManagedApps(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	result, err := listManagedApps(c)
	if err != nil {
		return err
	}
	if result == nil {
		result = []*ManagedApp{}
	}
	return c.JSON(200, result)
}

// CreateManagedApp implements ServerInterface.
func (s *server) CreateManagedApp(ec echo.Context) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	param := &ManagedAppParameters{}
	if err := c.Bind(param); err != nil {
		return err
	}

	result, err := createManagedApp(c, param)
	if err != nil {
		return err
	}

	return c.JSON(200, managedAppDocToModel(result))
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
