package agentadmin

import (
	"errors"
	"fmt"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// CreateAgent implements ServerInterface.
func (s *AgentAdminServer) CreateAgent(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	params := new(agentmodels.CreateAgentRequest)
	if err := c.Bind(params); err != nil {
		return err
	}

	if params.AppId != "" {
		// import existing application
		// check if application has owner
		return importAgentApp(c, params)
	} else if params.DisplayName != "" {
		return fmt.Errorf("%w: no display name set", base.ErrResponseStatusBadRequest)
	}

	c = c.Elevate()
	gclient := graph.GetServiceMsGraphClient(c)
	var application gmodels.Applicationable = gmodels.NewApplication()
	displayName := params.DisplayName
	application.SetDisplayName(&displayName)
	application.SetSignInAudience(to.Ptr("AzureADMyOrg"))
	apiApplication := gmodels.NewApiApplication()
	apiApplication.SetRequestedAccessTokenVersion(to.Ptr(int32(2)))
	application.SetApi(apiApplication)
	application, err := gclient.Applications().Post(c, application, nil)
	if err != nil {
		return err
	}
	doc := &AgentDoc{
		AppDoc: profile.AppDoc{
			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceProvider: models.NamespaceProviderProfile,
					NamespaceID:       profile.NamespaceIDApp,
					ResourceProvider:  models.ProfileResourceProviderAgent,
				},
				ID: *application.GetAppId(),
			},
			DisplayName:   application.GetDisplayName(),
			ApplicationID: application.GetId(),
		},
	}

	mSp := gmodels.NewServicePrincipal()
	mSp.SetAppId(application.GetAppId())
	if sp, err := gclient.ServicePrincipals().Post(c, mSp, nil); err != nil {
		return err
	} else {
		doc.ServicePrincipalID = sp.GetId()
		doc.ServicePrincipalType = sp.GetServicePrincipalType()
	}

	// persist document
	docService := resdoc.GetDocService(c)
	resp, err := docService.Create(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToProfile())

}

func importAgentApp(c ctx.RequestContext, params *agentmodels.CreateAgentRequest) error {
	systemappDoc, _, err := systemapp.GetSystemAppDoc(c, systemapp.SystemAppNameAPI)
	if err != nil {
		return err
	}
	if systemappDoc.ServicePrincipalID == nil || *systemappDoc.ServicePrincipalID == "" {
		return fmt.Errorf("%w: system application ID not found: %s, please sync", base.ErrResponseStatusNotFound, systemapp.SystemAppNameAPI)
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	application, err := gclient.ApplicationsWithAppId(&params.AppId).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "appId", "api", "signInAudience"},
		},
	})
	if err != nil {
		err = graph.HandleMsGraphError(err)
		if errors.Is(err, graph.ErrMsGraphResourceNotFound) {
			return fmt.Errorf("%w, application with AppId not found: %s", base.ErrResponseStatusNotFound, params.AppId)
		}
		return err
	}

	owners, err := gclient.Applications().ByApplicationId(*application.GetId()).Owners().Get(c, &applications.ItemOwnersRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ItemOwnersRequestBuilderGetQueryParameters{
			Select: []string{"id"},
		},
	})
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(owners.GetValue(), func(v gmodels.DirectoryObjectable) bool {
		log.Ctx(c).Debug().Interface("owner", v.GetId()).Msg("owner")

		return systemappDoc.ServicePrincipalID != nil && v.GetId() != nil && *systemappDoc.ServicePrincipalID == *v.GetId()
	}) {
		return fmt.Errorf("%w: application with AppId not owned by system application: %s", base.ErrResponseStatusForbidden, params.AppId)
	}

	agentDoc := &AgentDoc{
		AppDoc: profile.AppDoc{

			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceProvider: models.NamespaceProviderProfile,
					NamespaceID:       profile.NamespaceIDApp,
					ResourceProvider:  models.ProfileResourceProviderAgent,
				},
				ID: *application.GetAppId(),
			},
			DisplayName: application.GetDisplayName(),

			ApplicationID: application.GetId(),
		},
	}

	c = c.Elevate()
	gclient = graph.GetServiceMsGraphClient(c)

	sp, err := gclient.ServicePrincipalsWithAppId(application.GetAppId()).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "servicePrincipalType"},
		},
	})
	if err != nil {
		err = graph.HandleMsGraphError(err)
		if errors.Is(err, graph.ErrMsGraphResourceNotFound) {
			// create service principal
			mSp := gmodels.NewServicePrincipal()
			mSp.SetAppId(application.GetAppId())
			if sp, err = gclient.ServicePrincipals().Post(c, mSp, nil); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	agentDoc.ServicePrincipalID = sp.GetId()
	agentDoc.ServicePrincipalType = sp.GetServicePrincipalType()

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, agentDoc, nil)
	if err != nil {
		return err
	}

	// patch application if needed
	patchApplicaiton := gmodels.NewApplication()
	hasPatch := false

	if *application.GetSignInAudience() != "AzureADMyOrg" {
		patchApplicaiton.SetSignInAudience(to.Ptr("AzureADMyOrg"))
		hasPatch = true
	}
	if *application.GetApi().GetRequestedAccessTokenVersion() == 2 {
		appApi := application.GetApi()
		appApi.SetRequestedAccessTokenVersion(to.Ptr(int32(2)))
		patchApplicaiton.SetApi(appApi)
		hasPatch = true
	}
	if hasPatch {
		_, err = gclient.Applications().ByApplicationId(*application.GetId()).Patch(c, patchApplicaiton, nil)
		if err != nil {
			return err
		}
	}

	return c.JSON(resp.RawResponse.StatusCode, agentDoc.ToProfile())
}
