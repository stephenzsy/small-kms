package profile

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

type ProfileContext interface {
	GetResourceDocNsID() kmsdoc.DocNsID
	GetSelfProfileDoc(c common.ServiceContext) (*ProfileDoc, error)
	GetProfileDoc(common.ServiceContext, kmsdoc.DocNsID, kmsdoc.DocID) (*ProfileDoc, error)
	GetRequestProfileType() models.ProfileType
}

type profileContext struct {
	service            *profileService
	profileNsID        kmsdoc.DocNsID
	profileID          kmsdoc.DocID
	requestProfileType models.ProfileType
}

// GetRequestProfileType implements ProfileContextService.
func (p *profileContext) GetRequestProfileType() models.ProfileType {
	return p.requestProfileType
}

// GetProfileDoc implements ProfileContextService.
func (p *profileContext) GetProfileDoc(c common.ServiceContext, nsID kmsdoc.DocNsID, docID kmsdoc.DocID) (*ProfileDoc, error) {
	return getProfileDoc(c, nsID, docID)
}

func GetResourceNsIDForProfile(profileID kmsdoc.DocID) kmsdoc.DocNsID {
	switch profileID.Kind() {
	case kmsdoc.DocKindCaRoot:
		return kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaRoot, profileID.Identifier())
	case kmsdoc.DocKindCaInt:
		return kmsdoc.NewDocIdentifier(kmsdoc.DocNsTypeCaInt, profileID.Identifier())
	}
	return kmsdoc.NewDocIdentifier(kmsdoc.DocNSTypeDirectory, profileID.Identifier())
}

// GetResourceDocNsID implements ProfileContextService.
func (p *profileContext) GetResourceDocNsID() kmsdoc.DocNsID {
	return GetResourceNsIDForProfile(p.profileID)
}

// GetSelfProfileDoc implements ProfileContextService.
func (p *profileContext) GetSelfProfileDoc(c common.ServiceContext) (*ProfileDoc, error) {
	return getProfileDoc(c, p.profileNsID, p.profileID)
}

func (s *profileService) newProfileContext(profileType models.ProfileType, identifier common.Identifier) (c profileContext, err error) {
	c.service = s
	c.profileNsID, c.profileID, err = GetProfileInternalIDs(profileType, identifier)
	c.requestProfileType = profileType
	return
}

var _ ProfileContext = (*profileContext)(nil)
