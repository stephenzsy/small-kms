package managedapp

import (
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type server struct {
}

// ListManagedApps implements ServerInterface.
func (*server) ListManagedApps(ec echo.Context) error {
	c := ctx.ResolveRequestContext(ec)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return err
	}

	result, err := listManagedApps(c)
	if err != nil {
		return err
	}
	if result == nil {
		result = []*ManagedAppDoc{}
	}
	return c.JSON(200, result)
}

// CreateManagedApp implements ServerInterface.
func (*server) CreateManagedApp(ec echo.Context) error {
	c := ctx.ResolveRequestContext(ec)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return err
	}

	param := &ManagedAppParameters{}
	if err := c.Bind(param); err != nil {
		return err
	}

	result, err := createManagedApp(c, param)
	if err != nil {
		return err
	}

	return c.JSON(200, result)
}

var _ ServerInterface = (*server)(nil)

func NewServer() *server {
	return &server{}
}
