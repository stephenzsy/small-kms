package managedapp

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
)

// SyncManagedApp implements ServerInterface.
func (s *server) SyncManagedApp(ec echo.Context, managedAppId uuid.UUID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	application, err := gclient.ApplicationsWithAppId(to.Ptr(managedAppId.String())).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
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

	// query service principal
	sp, err := gclient.ServicePrincipalsWithAppId(to.Ptr(managedAppId.String())).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "servicePrincipalType"},
		},
	})
	if err != nil {
		return err
	}

	doc := &ManagedAppDoc{}
	doc.Init(managedAppId, *application.GetDisplayName(), namespaceIDNameManagedApp)
	if doc.ServicePrincipalID, err = uuid.Parse(*sp.GetId()); err != nil {
		return err
	}
	doc.ServicePrincipalType = *sp.GetServicePrincipalType()
	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Upsert(c, doc, nil); err != nil {
		return err
	}
	m := new(ManagedApp)
	doc.PopulateModel(m)
	return c.JSON(200, m)
}
