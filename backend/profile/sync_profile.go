package profile

import (
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// SyncProfile implements ProfileService.
func SyncProfile(c common.ServiceContext) (*models.ProfileComposed, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	identifier := nsID.Identifier()

	if id, ok := identifier.TryGetUUID(); !ok || id.Version() != 4 {
		return nil, fmt.Errorf("%w:invalid profile id for sync", common.ErrStatusBadRequest)
	}

	client, err := common.GetClientProvider(c).MsGraphDelegatedClient(c)
	if err != nil {
		return nil, err
	}
	directoryObjId := identifier.String()
	dirObject, err := client.DirectoryObjects().ByDirectoryObjectId(directoryObjId).Get(c, nil)
	profileLocator := resolveProfileLocatorFromNamespaceID(nsID)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("directoryObject:%s", directoryObjId))
		if errors.Is(err, common.ErrStatusNotFound) {
			// delete existing profile if exists
			err = kmsdoc.DeleteByRef(c, profileLocator)
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
