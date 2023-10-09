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
	"github.com/stephenzsy/small-kms/backend/utils"
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
	}
	if doc == nil {
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
	} else {
		ops.AppendSet("/graphSyncCode", "")
		ops.AppendSet("/graph", profileDoc.Graph)
		ops.AppendSet("/@odata.type", profileDoc.OdataType)
		ops.AppendSet("/displayName", profileDoc.DispalyName)
	}
	err = kmsdoc.Patch(c, profileDoc.GetLocator(), profileDoc, ops, &azcosmos.ItemOptions{
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

	client, err := c.ServiceClientProvider().MsGraphDelegatedClient(c)
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

func createManagedApplication(c RequestContext, req models.CreateManagedApplicationProfileRequest) (*Profile, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w:invalid name", common.ErrStatusBadRequest)
	}
	app := msgraphmodels.NewApplication()
	app.SetDisplayName(&req.Name)
	app.SetSignInAudience(utils.ToPtr("AzureADMyOrg"))
	client := c.ServiceClientProvider().MsGraphServerClient()
	applicationable, err := client.Applications().Post(c, app, nil)
	if err != nil {
		return nil, err
	}

	profileDoc := ProfileDoc{}
	err = profileDoc.init(applicationable)
	if err != nil {
		return nil, err
	}
	profileDoc.IsAppManaged = utils.ToPtr(true)
	err = kmsdoc.Upsert(c, &profileDoc)
	return profileDoc.toModel(), err
}

func createApplicationManagedServicePrincipal(c RequestContext, ownerDoc *ProfileDoc) (*Profile, error) {
	delegatedClient, err := c.ServiceClientProvider().MsGraphDelegatedClient(c)
	if err != nil {
		return nil, err
	}
	if ownerDoc.GraphSyncCode != "" {
		return nil, fmt.Errorf("%w:invalid application profile for creating managed service principal, status: %s",
			common.ErrStatusBadRequest, ownerDoc.GraphSyncCode)
	}
	if ownerDoc.Graph == nil || ownerDoc.Graph.AppID == nil || *ownerDoc.Graph.AppID == "" {
		return nil, fmt.Errorf("%w:invalid application profile for creating managed service principal, missing appId", common.ErrStatusBadRequest)
	}
	sp, err := delegatedClient.ServicePrincipalsWithAppId(ownerDoc.Graph.AppID).Get(c, nil)
	if err != nil {
		err = common.WrapMsGraphNotFoundErr(err, "service principal")
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, err
		}
	}
	if sp == nil {
		// create service principal
		mSp := msgraphmodels.NewServicePrincipal()
		mSp.SetAppId(ownerDoc.Graph.AppID)
		client := c.ServiceClientProvider().MsGraphServerClient()
		sp, err = client.ServicePrincipals().Post(c, mSp, nil)
		if err != nil {
			return nil, err
		}
	}
	profileDoc := ProfileDoc{}
	err = profileDoc.init(sp)
	if err != nil {
		return nil, err
	}
	profileDoc.IsAppManaged = utils.ToPtr(true)
	profileDoc.Owner = ownerDoc.GetLocator()
	cElevated := c.Elevate() // elevate to avoid cancellation
	err = kmsdoc.Upsert(cElevated, &profileDoc)
	if err != nil {
		return nil, err
	}
	patchOps := azcosmos.PatchOperations{}
	if ownerDoc.Owns == nil {
		patchOps.AppendSet("/@owns", map[NamespaceKind]models.ResourceLocator{
			models.NamespaceKindServicePrincipal: profileDoc.GetLocator(),
		})
	} else {
		patchOps.AppendSet(fmt.Sprintf("/@owns/%s", models.NamespaceKindServicePrincipal), profileDoc.GetLocator())
	}
	err = kmsdoc.Patch(cElevated, ownerDoc.GetLocator(), ownerDoc, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: &ownerDoc.ETag,
	})
	return profileDoc.toModel(), err
}

func CreateProfile(c RequestContext, namespaceKind models.NamespaceKind, req CreateProfileRequest) (*Profile, error) {
	// validate name
	if req, err := req.AsCreateManagedApplicationProfileRequest(); err == nil {
		if namespaceKind != models.NamespaceKindApplication {
			return nil, fmt.Errorf("%w:invalid namespace kind for creating managed application profile", common.ErrStatusBadRequest)
		}
		return createManagedApplication(c, req)
	}
	discriminiator, _ := req.Discriminator()
	return nil, fmt.Errorf("%w:bad request type: %s", common.ErrStatusBadRequest, discriminiator)
}

func CreateManagedProfile(c RequestContext, targetNamespaceKind models.NamespaceKind) (*Profile, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	ownerDoc, err := GetResourceProfileDoc(c)
	if err != nil {
		return nil, err
	}
	switch targetNamespaceKind {
	case models.NamespaceKindServicePrincipal:
		switch nsID.Kind() {
		case models.NamespaceKindApplication:
			return createApplicationManagedServicePrincipal(c, ownerDoc)
		}
	}
	return nil, fmt.Errorf("%w:invalid target namespace kind: %s", common.ErrStatusBadRequest, targetNamespaceKind)
}
