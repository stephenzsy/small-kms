package profile

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

func getProfileDoc(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (doc *ProfileDoc, err error) {
	if profileType == models.ProfileTypeRootCA {
		if a, ok := rootCaProfileDocs[identifier]; ok {
			return &a, nil
		}
		return nil, common.ErrStatusNotFound
	}
	if profileType == models.ProfileTypeIntermediateCA {
		if a, ok := rootCaProfileDocs[identifier]; ok {
			return &a, nil
		}
		return nil, common.ErrStatusNotFound
	}

	doc = &ProfileDoc{}
	err = kmsdoc.Read(c, docNsIDProfileTenant, kmsdoc.NewDocIdentifier(kmsdoc.DocTypeDirectoryObject, identifier), doc)
	return
}

// GetProfile implements ProfileService.
func (*profileService) GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}
	doc, err := getProfileDoc(c, profileType, identifier)
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}
