package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/models"
)

// GetServiceConfig implements models.ServerInterface.
func (*server) GetServiceConfig(ctx echo.Context) error {
	bad := func(e error) error {
		return wrapResponse[*models.ServiceConfigComposed](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)

	}

	result, err := admin.GetServiceConfig(c)
	return wrapResponse[*models.ServiceConfigComposed](c, http.StatusOK, result, err)
}
