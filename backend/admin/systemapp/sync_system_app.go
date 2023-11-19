package systemapp

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// SyncSystemApp implements ServerInterface.
func (s *SystemAppAdminServer) SyncSystemApp(ec echo.Context, id string) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	appName, err := validateSystemAppName(id)
	if err != nil {
		return err
	}

	appID, err := resolveSystemAppID(c, appName)
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
	doc := &SystemAppDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderProfile,
				NamespaceID:       profile.NamespaceIDApp,
				ResourceProvider:  models.ProfileResourceProviderSystem,
			},
			ID: appID.String(),
		},
		DisplayName:          sp.GetDisplayName(),
		ServicePrincipalID:   sp.GetId(),
		ServicePrincipalType: sp.GetServicePrincipalType(),
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
		doc.ApplicationID = application.GetId()
	}

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToProfile())

}
