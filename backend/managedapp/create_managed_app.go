package managedapp

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
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
	apiApplication := gmodels.NewApiApplication()
	apiApplication.SetRequestedAccessTokenVersion(to.Ptr(int32(2)))
	application.SetApi(apiApplication)
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

func apiSyncManagedApp(c ctx.RequestContext, appID uuid.UUID) error {
	gclient := graph.GetServiceMsGraphClient(c)
	application, err := gclient.ApplicationsWithAppId(to.Ptr(appID.String())).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "api", "displayName", "appId"},
		},
	})
	if err != nil {
		return err
	}

	var patchApplication *gmodels.Application

	if application.GetApi().GetRequestedAccessTokenVersion() == nil || *application.GetApi().GetRequestedAccessTokenVersion() != 2 {
		if patchApplication == nil {
			patchApplication = gmodels.NewApplication()
		}
		if patchApplication.GetApi() == nil {
			patchApplication.SetApi(gmodels.NewApiApplication())
		}
		patchApplication.GetApi().SetRequestedAccessTokenVersion(to.Ptr(int32(2)))
	}

	if patchApplication != nil {
		_, err = gclient.Applications().ByApplicationId(*application.GetId()).Patch(c, patchApplication, nil)
		if err != nil {
			return err
		}
	}

	// serviceAppA, serviceAppSP, err := graph.GetServiceAppAndSP(c, gclient)
	// if err != nil {
	// 	return err
	// }
	// for _, ar := range serviceAppSP.GetAppRoles() {
	// 	log.Debug().Any("id", ar.GetId()).Any("value", ar.GetValue()).Msg("serviceApplication")
	// }

	doc, err := getManagedApp(c, appID)
	if err != nil {
		return err
	}
	m := new(ManagedApp)
	doc.PopulateModel(m)
	return c.JSON(200, m)
}
