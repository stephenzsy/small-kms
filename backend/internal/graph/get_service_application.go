package graph

import (
	"context"
	"errors"

	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
)

// should use service client
func GetServiceAppAndSP(c context.Context, client *msgraph.GraphServiceClient) (gmodels.Applicationable, gmodels.ServicePrincipalable, error) {
	if appID, ok := c.Value(ServiceMsGraphClientClientIDContextKey).(string); ok {
		app, err := client.ApplicationsWithAppId(&appID).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id", "appRoles"},
			},
		})
		if err != nil {
			return nil, nil, err
		}

		sp, err := client.ServicePrincipalsWithAppId(&appID).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id", "appRoles"},
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return app, sp, nil
	}
	return nil, nil, errors.New("no app id in context")
}
