package profile

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

func getProfileDoc(c common.ServiceContext, nsID kmsdoc.DocNsID, docID kmsdoc.DocID) (doc *ProfileDoc, err error) {
	if nsID == docNsIDProfileBuiltIn {
		if docID.Kind() == kmsdoc.DocKindCaRoot {
			if a, ok := rootCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		}
		if docID.Kind() == kmsdoc.DocKindCaInt {
			if a, ok := rootCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		}
	}
	doc = &ProfileDoc{}
	err = kmsdoc.Read(c, docNsIDProfileTenant, docID, doc)
	return
}

// GetProfile implements ProfileService.
func (*profileService) GetProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}
	nsID, docID, err := GetProfileInternalIDs(profileType, identifier)
	if err != nil {
		return nil, err
	}
	doc, err := getProfileDoc(c, nsID, docID)
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}
