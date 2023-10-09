package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
)

// ListProfiles implements models.ServerInterface.
func (s *server) ListProfiles(ec echo.Context, namespaceKind models.NamespaceKind) error {
	c := ec.(RequestContext)
	respData, respErr := (func() ([]*models.ProfileRefComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		return profile.ListProfiles(c, namespaceKind)
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
func (s *server) SyncProfile(ec echo.Context, namespaceKind models.NamespaceKind, namespaceId models.Identifier) error {
	c := ec.(RequestContext)
	respData, respErr := (func() (*models.ProfileComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
		if err != nil {
			return nil, err
		}
		return profile.SyncProfile(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// CreateProfile implements models.ServerInterface.
func (*server) CreateProfile(ctx echo.Context, namespaceKind models.NamespaceKind) error {
	bad := func(e error) error {
		return wrapResponse[*models.ProfileComposed](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)
	}

	req := models.CreateProfileRequest{}
	if err := c.Bind(&req); err != nil {
		return bad(err)
	}

	result, err := profile.CreateProfile(c, namespaceKind, req)
	return wrapResponse[*models.ProfileComposed](c, http.StatusOK, result, err)
}

// CreateManagedNamespace implements models.ServerInterface.
func (*server) CreateManagedNamespace(ctx echo.Context,
	namespaceKind models.NamespaceKind,
	namespaceId common.Identifier,
	targetNamespaceKind models.NamespaceKind) error {
	bad := func(e error) error {
		return wrapResponse[*models.ProfileComposed](ctx, http.StatusOK, nil, e)
	}
	c := ctx.(RequestContext)

	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return bad(err)

	}

	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	result, err := profile.CreateManagedProfile(c, targetNamespaceKind)
	return wrapResponse[*models.ProfileComposed](c, http.StatusOK, result, err)
}
