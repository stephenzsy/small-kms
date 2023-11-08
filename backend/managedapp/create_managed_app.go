package managedapp

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
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
		doc.Init(appID, displayName, namespaceIDNameManagedApp)
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

func resolveSystemAppID(c context.Context, systemAppName SystemAppName) (uuid.UUID, error) {
	switch systemAppName {
	case SystemAppNameBackend:
		if systemAppID, ok := c.Value(graph.ServiceClientIDContextKey).(string); ok {
			return uuid.Parse(systemAppID)
		}
	case SystemAppNameAPI:
		if systemAppID, ok := c.Value(graph.ServiceMsGraphClientClientIDContextKey).(string); ok {
			return uuid.Parse(systemAppID)
		}
	}
	return uuid.Nil, fmt.Errorf("%w: system app not found: %s", base.ErrResponseStatusNotFound, systemAppName)
}

func apiSyncSystemApp(c ctx.RequestContext, systemAppName SystemAppName) error {
	appID, err := resolveSystemAppID(c, systemAppName)
	if err != nil {
		return err
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	sp, err := gclient.ServicePrincipalsWithAppId(to.Ptr(appID.String())).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "appId", "servicePrincipalType"},
		},
	})
	if err != nil {
		return err
	}
	doc := &ManagedAppDoc{}
	doc.Init(appID, *sp.GetDisplayName(), namespaceIDNameSystemApp)

	if doc.ServicePrincipalID, err = uuid.Parse(*sp.GetId()); err != nil {
		return err
	}
	if sp.GetServicePrincipalType() != nil {
		doc.ServicePrincipalType = *sp.GetServicePrincipalType()
	}

	if *sp.GetServicePrincipalType() == "Application" {
		application, err := gclient.ApplicationsWithAppId(to.Ptr(appID.String())).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "appId"},
			},
		})
		if err != nil {
			return err
		}
		if doc.ApplicationID, err = uuid.Parse(*application.GetId()); err != nil {
			return err
		}
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	err = docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	m := new(ManagedApp)
	doc.PopulateModel(m)
	return c.JSON(200, m)
}
