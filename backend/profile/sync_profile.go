package profile

import (
	"errors"
	"fmt"

	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func StoreProfile(c RequestContext, dirObject msgraphmodels.DirectoryObjectable) (*ProfileDoc, error) {
	profileDoc := ProfileDoc{}
	err := profileDoc.init(dirObject)
	if err != nil {
		return nil, err
	}

	err = kmsdoc.Upsert(c, &profileDoc)
	if err != nil {
		return nil, err
	}

	return &profileDoc, nil
}

// SyncProfile implements ProfileService.
func SyncProfile(c RequestContext) (*models.ProfileComposed, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	identifier := nsID.Identifier()

	if id, ok := identifier.TryGetUUID(); !ok || id.Version() != 4 {
		return nil, fmt.Errorf("%w:invalid profile id for sync", common.ErrStatusBadRequest)
	}

	client, err := c.ServiceClientProvider().MsGraphDelegatedClient(c)
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
	pdoc, err := StoreProfile(c, dirObject)
	if err != nil {
		return nil, err
	}
	return pdoc.toModel(), nil
}
