package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/models"
)

// GetServiceConfig implements models.ServerInterface.
func (*server) GetServiceConfig(ctx echo.Context) error {

	c := ctx.(RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return respondRequireAdmin(c)
	}

	result, err := admin.GetServiceConfig(c)
	return wrapResponse[*models.ServiceConfigComposed](c, http.StatusOK, result, err)
}

// PatchServiceConfig implements models.ServerInterface.
func (*server) PatchServiceConfig(ctx echo.Context, configPath models.PatchServiceConfigParamsConfigPath) error {

	c := ctx.(RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return respondRequireAdmin(c)
	}

	result, err := admin.PatchServiceConfig(c, configPath)
	return wrapResponse[*models.ServiceConfigComposed](c, http.StatusOK, result, err)
}
