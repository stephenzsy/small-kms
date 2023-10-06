package profile

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

// GetProfile implements ProfileService.
func (*profileService) GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}
	if profileType == models.ProfileTypeRootCA {
		if a, ok := rootCaProfiles[identifier]; ok {
			return &a, nil
		}
		return nil, common.ErrStatusNotFound
	}
	if profileType == models.ProfileTypeIntermediateCA {
		if a, ok := intCaProfiles[identifier]; ok {
			return &a, nil
		}
		return nil, common.ErrStatusNotFound
	}

	doc := ProfileDoc{}
	err := kmsdoc.ReadByKeyFunc(c, func() (string, string) {
		return getProfileDocKey(identifier)
	}, &doc)
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}
