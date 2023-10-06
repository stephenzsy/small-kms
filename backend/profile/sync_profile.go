package profile

import (
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
)

// SyncProfile implements ProfileService.
func (s *profileService) SyncProfile(c common.ServiceContext, profileType models.ProfileType, identifier models.Identifier) (*models.Profile, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	if id, ok := identifier.TryGetUUID(); !ok || id.Version() != 4 {
		return nil, fmt.Errorf("%w:invalid profile id", common.ErrStatusBadRequest)
	}

	client, err := common.GetClientProvider(c).MsGraphDelegatedClient(c)
	if err != nil {
		return nil, err
	}
	directoryObjId := identifier.String()
	dirObject, err := client.DirectoryObjects().ByDirectoryObjectId(directoryObjId).Get(c, nil)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("directoryObject:%s", directoryObjId))
		if errors.Is(err, common.ErrStatusNotFound) {
			// delete existing profile if exists
			err = kmsdoc.DeleteByKey(c, docNsIDProfileTenant, kmsdoc.NewDocIdentifier(kmsdoc.DocKindDirectoryObject, identifier))
		}
		return nil, err
	}
	profileDoc := ProfileDoc{}
	err = profileDoc.init(dirObject)
	if err != nil {
		return nil, err
	}

	err = kmsdoc.Upsert(c, &profileDoc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
