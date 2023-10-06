package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephenzsy/small-kms/backend/models"
)

// GetProfile implements models.ServerInterface.
func (s *server) GetProfile(c *gin.Context, profileType models.ProfileType, identifier models.Identifier) {
	res, err := s.profileService.GetProfile(s.ServiceContext(c), profileType, identifier)
	wrapResponse(c, http.StatusOK, res, err)
}

// ListProfiles implements models.ServerInterface.
func (s *server) ListProfiles(c *gin.Context, profileType models.ProfileType) {
	res, err := s.profileService.ListProfiles(s.ServiceContext(c), profileType)
	wrapResponse(c, http.StatusOK, res, err)
}

// SyncProfile implements models.ServerInterface.
func (s *server) SyncProfile(c *gin.Context, profileType models.ProfileType, identifier models.Identifier) {
	res, err := s.profileService.SyncProfile(s.ServiceContext(c), profileType, identifier)
	wrapResponse(c, http.StatusOK, res, err)
}
