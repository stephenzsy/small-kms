package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
)

// ListProfiles implements models.ServerInterface.
func (s *server) ListProfiles(c *gin.Context, params models.ListProfilesParams) {
	respData, respErr := (func() ([]*models.ProfileRefComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}
		sc := s.ServiceContext(c)

		return profile.ListProfiles(sc, params.ProfileType)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}

// GetProfile implements models.ServerInterface.
func (s *server) GetProfile(c *gin.Context, profileType models.NamespaceKind, identifier models.Identifier) {
	respData, respErr := (func() (*models.ProfileComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}
		sc := s.ServiceContext(c)
		sc, err := ns.WithNamespaceContext(sc, profileType, identifier)
		if err != nil {
			return nil, err
		}
		return profile.GetProfile(sc)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)
}

// SyncProfile implements models.ServerInterface.
func (s *server) SyncProfile(c *gin.Context, profileType models.NamespaceKind, identifier models.Identifier) {
	respData, respErr := (func() (*models.ProfileComposed, error) {

		if err := auth.AuthorizeAdminOnly(c); err != nil {
			return nil, err
		}

		sc := s.ServiceContext(c)
		sc, err := ns.WithNamespaceContext(sc, profileType, identifier)
		if err != nil {
			return nil, err
		}
		return profile.SyncProfile(sc)
	})()
	wrapResponse(c, http.StatusOK, respData, respErr)

}
