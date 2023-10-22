package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// ListProfiles implements models.ServerInterface.
func (s *server) ListProfiles(ec echo.Context, namespaceKind models.NamespaceKind) error {
	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	respData, respErr := profile.ListProfiles(c, namespaceKind)
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// GetProfile implements models.ServerInterface.
func (s *server) GetProfile(ec echo.Context, profileType models.NamespaceKind, identifier shared.Identifier) error {
	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	respData, respErr := (func() (*models.ProfileComposed, error) {
		c, err := ns.WithNamespaceContext(c, profileType, identifier)
		if err != nil {
			return nil, err
		}
		return profile.GetProfile(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// SyncProfile implements models.ServerInterface.
func (s *server) SyncProfile(ec echo.Context, namespaceKind models.NamespaceKind, namespaceId shared.Identifier) error {
	c := ec.(RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return respondRequireAdmin(c)
	}
	respData, respErr := (func() (*models.ProfileComposed, error) {

		c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
		if err != nil {
			return nil, err
		}
		return profile.SyncProfile(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}
