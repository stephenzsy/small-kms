package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type ProfileService interface {
	GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	SyncProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	ListProfiles(c common.ServiceContext, profileType models.ProfileType) ([]*models.ProfileRef, error)
	WithProfileContext(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (common.ServiceContext, error)
}

type profileService struct {
}

func NewProfileService() ProfileService {
	return &profileService{}
}

type profileContextKeyType string

const profileContextKey profileContextKeyType = "profileContext"

func (s *profileService) WithProfileContext(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (common.ServiceContext, error) {
	pcs, err := s.newProfileContext(profileType, identifier)
	return context.WithValue(c, profileContextKey, &pcs), err
}

func GetProfileContextService(c common.ServiceContext) ProfileContextService {
	if pc, ok := c.Value(profileContextKey).(ProfileContextService); ok {
		return pc
	}
	return nil
}
