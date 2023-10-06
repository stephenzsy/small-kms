package profile

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type ProfileService interface {
	GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	SyncProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	ListProfiles(c common.ServiceContext, profileType models.ProfileType) ([]models.ProfileRef, error)
}

type profileService struct {
}

func NewProfileService() ProfileService {
	return &profileService{}
}
