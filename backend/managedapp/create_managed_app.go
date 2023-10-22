package managedapp

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
)

func createApplicationManagedServicePrincipal(c context.Context, gclient *msgraphsdkgo.GraphServiceClient, docService base.AzCosmosCRUDDocService, appDoc *ManagedAppDoc) error {

	// create service principal
	mSp := gmodels.NewServicePrincipal()
	mSp.SetAppId(to.Ptr[string](appDoc.GetAppID().String()))
	sp, err := gclient.ServicePrincipals().Post(c, mSp, nil)
	if err != nil {
		return err
	}

	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet(patchColumnServicePrincipalID, sp.GetId())
	appDoc.ServicePrincipalID, err = uuid.Parse(*sp.GetId())
	if err != nil {
		return err
	}
	err = docService.Patch(c, appDoc, patchOps, nil)
	if err != nil {
		return err
	}
	return nil
}

func createManagedApp(c context.Context, params *ManagedAppParameters) (*ManagedAppDoc, error) {
	bad := func(e error) (*ManagedAppDoc, error) {
		return nil, e
	}

	c = ctx.Elevate(c)
	gclient := graph.GetServiceMsGraphClient(c)
	application := gmodels.NewApplication()
	displayName := params.DisplayName
	application.SetDisplayName(&displayName)
	application.SetSignInAudience(to.Ptr("AzureADMyOrg"))
	if application, err := gclient.Applications().Post(c, application, nil); err != nil {
		return bad(err)
	} else {
		appID, err := uuid.Parse(*application.GetAppId())
		if err != nil {
			return bad(err)
		}
		doc := &ManagedAppDoc{}
		doc.Init(appID, displayName)
		doc.ApplicationID, err = uuid.Parse(*application.GetId())
		if err != nil {
			return bad(err)
		}

		// persist document
		docService := base.GetAzCosmosCRUDService(c)
		err = docService.Create(c, doc, nil)
		if err != nil {
			return bad(err)
		}

		if params.SkipServicePrincipalCreation == nil || !*params.SkipServicePrincipalCreation {
			err = createApplicationManagedServicePrincipal(c, gclient, docService, doc)
			if err != nil {
				return bad(err)
			}
		}
		return doc, nil
	}
}
