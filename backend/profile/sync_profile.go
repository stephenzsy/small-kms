package profile

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func StoreProfile(c RequestContext, dirObject msgraphmodels.DirectoryObjectable, odataErrorCode *string, graphErr error) (*ProfileDoc, error) {
	profileDoc := ProfileDoc{}
	err := profileDoc.init(dirObject)
	if err != nil {
		return nil, err
	}
	return upsertProfileDoc(c, &profileDoc, odataErrorCode, graphErr)
}

func upsertProfileDoc(c RequestContext, profileDoc *ProfileDoc, odataErrorCode *string, graphErr error) (*ProfileDoc, error) {
	// load existing profile
	doc, err := getProfileDoc(c, profileDoc.GetLocator())
	if err != nil {
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, err
		}

		// no existing doc, create new
		if graphErr != nil {
			return nil, graphErr
		}
		err = kmsdoc.Create(c, profileDoc)
		return profileDoc, err
	}
	// has existing doc, patch
	ops := azcosmos.PatchOperations{}
	if graphErr != nil {
		ops.AppendSet("/graphSyncCode", odataErrorCode)
		profileDoc.GraphSyncCode = *odataErrorCode
	} else {
		ops.AppendSet("/graphSyncCode", "")
		ops.AppendSet("/graph", profileDoc.Graph)
		ops.AppendSet("/@odata.type", profileDoc.OdataType)
		ops.AppendSet("/displayName", profileDoc.DispalyName)
	}
	err = kmsdoc.Patch(c, profileDoc, ops, &azcosmos.ItemOptions{
		IfMatchEtag: &doc.ETag,
	})
	if err != nil {
		return nil, err
	}

	return profileDoc, graphErr
}

// SyncProfile implements ProfileService.
func SyncProfile(c RequestContext) (*models.ProfileComposed, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	identifier := nsID.Identifier()

	if id, ok := identifier.TryGetUUID(); !ok || id.Version() != 4 {
		return nil, fmt.Errorf("%w:invalid profile id for sync", common.ErrStatusBadRequest)
	}

	client, err := common.GetAdminServerRequestClientProvider(c).MsGraphClient()
	if err != nil {
		return nil, err
	}
	directoryObjId := identifier.String()
	var getGraphErrorCode *string
	dirObject, err := client.DirectoryObjects().ByDirectoryObjectId(directoryObjId).Get(c, nil)
	if err != nil {
		var isODataError bool
		if getGraphErrorCode, _, isODataError = common.ExtractGraphODataErrorCode(err); !isODataError {
			return nil, err
		}
	}
	pdoc, err := StoreProfile(c, dirObject, getGraphErrorCode, err)
	if err != nil {
		return nil, err
	}
	return pdoc.toModel(), nil
}
