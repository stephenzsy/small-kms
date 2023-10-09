package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
)

// ListProfiles implements models.ServerInterface.
func (s *server) ListProfiles(ec echo.Context, params models.ListProfilesParams) error {
	c := ec.(RequestContext)
	respData, respErr := (func() ([]*models.ProfileRefComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		return profile.ListProfiles(c, params.ProfileType)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// GetProfile implements models.ServerInterface.
func (s *server) GetProfile(ec echo.Context, profileType models.NamespaceKind, identifier models.Identifier) error {
	c := ec.(RequestContext)
	respData, respErr := (func() (*models.ProfileComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}
		c, err := ns.WithNamespaceContext(c, profileType, identifier)
		if err != nil {
			return nil, err
		}
		return profile.GetProfile(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// SyncProfile implements models.ServerInterface.
func (s *server) SyncProfile(ec echo.Context, profileType models.NamespaceKind, identifier models.Identifier) error {
	c := ec.(RequestContext)
	respData, respErr := (func() (*models.ProfileComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		c, err := ns.WithNamespaceContext(c, profileType, identifier)
		if err != nil {
			return nil, err
		}
		return profile.SyncProfile(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}
