package profile

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type ProfileService interface {
	GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	SyncProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error)
	ListProfiles(c common.ServiceContext, profileType models.ProfileType) ([]models.ProfileRef, error)
	WithProfileContext(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) common.ServiceContext
}

type profileService struct {
}

func NewProfileService() ProfileService {
	return &profileService{}
}

type ProfileContext interface {
	ProfileType() models.ProfileType
	Identifier() models.Identifier
	EnsureProfileEnabled(c common.ServiceContext) error
}

type profileContext struct {
	service     *profileService
	profileType models.ProfileType
	identifier  models.Identifier
}

// Identifier implements ProfileContext.
func (c *profileContext) Identifier() models.NameOrUUIDIdentifier {
	return c.identifier
}

// ProfileType implements ProfileContext.
func (c *profileContext) ProfileType() models.ProfileType {
	return c.profileType
}

// Service implements ProfileContext.
func (pc *profileContext) EnsureProfileEnabled(c common.ServiceContext) error {
	profile, err := pc.service.GetProfile(c, pc.profileType, pc.identifier)
	if err != nil {
		return err
	}
	if profile.Metadata.Deleted != nil && !profile.Metadata.Deleted.IsZero() {
		return fmt.Errorf("%w:profile deleted", common.ErrStatusBadRequest)
	}
	return nil
}

type profileContextKeyType string

const profileContextKey profileContextKeyType = "profileContext"

func (s *profileService) WithProfileContext(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) common.ServiceContext {
	var pc ProfileContext = &profileContext{
		service:     s,
		profileType: profileType,
		identifier:  identifier,
	}
	return context.WithValue(c, profileContextKey, pc)
}

func GetProfileContext(c common.ServiceContext) ProfileContext {
	if pc, ok := c.Value(profileContextKey).(ProfileContext); ok {
		return pc
	}
	return nil
}
