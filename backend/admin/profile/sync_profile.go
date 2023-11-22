package profile

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"

	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// SyncProfile implements admin.ServerInterface.
func (*ProfileServer) SyncProfile(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string) error {

	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	switch namespaceProvider {
	case models.NamespaceProviderServicePrincipal,
		models.NamespaceProviderGroup,
		models.NamespaceProviderUser:
		// ok
	default:
		return base.ErrResponseStatusBadRequest
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}

	doc, err := SyncProfileInternal(c, namespaceId, gclient)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.ToModel())
}

func SyncProfileInternal(c ctx.RequestContext, namespaceId string, gclient *msgraphsdkgo.GraphServiceClient) (*ProfileDoc, error) {
	bad := func(e error) (*ProfileDoc, error) {
		return nil, e
	}

	doc := &ProfileDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderProfile,
				NamespaceID:       NamespaceIDGraph,
			},
			ID: namespaceId,
		},
	}
	dirObj, err := gclient.DirectoryObjects().ByDirectoryObjectId(namespaceId).Get(c, &directoryobjects.DirectoryObjectItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &directoryobjects.DirectoryObjectItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "appId", "servicePrincipalType", "userPrincipalName", "mail"},
		},
	})
	if err != nil {
		err = graph.HandleMsGraphError(err)
		if errors.Is(err, graph.ErrMsGraphResourceNotFound) {
			return bad(fmt.Errorf("%w,%w", base.ErrResponseStatusNotFound, err))
		}
		return bad(err)
	}
	switch *dirObj.GetOdataType() {
	case "#microsoft.graph.servicePrincipal":
		sp := dirObj.(gmodels.ServicePrincipalable)
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderServicePrincipal
		doc.DisplayName = sp.GetDisplayName()
		doc.ID = *sp.GetId()
		doc.AppId = sp.GetAppId()
		doc.ServicePrincipalType = sp.GetServicePrincipalType()
	case "#microsoft.graph.group":
		grp := dirObj.(gmodels.Groupable)
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderGroup
		doc.DisplayName = grp.GetDisplayName()
		doc.ID = *grp.GetId()
	case "#microsoft.graph.user":
		usr := dirObj.(gmodels.Userable)
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderUser
		doc.DisplayName = usr.GetDisplayName()
		doc.ID = *usr.GetId()
		doc.UserPrincipalName = usr.GetUserPrincipalName()
		doc.Mail = usr.GetMail()
		// ok
	default:
		return bad(fmt.Errorf("%w: object type is not supported %s not supported", base.ErrResponseStatusBadRequest, *dirObj.GetOdataType()))
	}

	_, err = resdoc.GetDocService(c).Upsert(c, doc, nil)
	return doc, err
}

func SyncServicePrincipalProfile(c ctx.RequestContext, ID string, additionalSelects []string) (ctx.RequestContext, *ProfileDoc, gmodels.ServicePrincipalable, error) {
	return syncServicePrincipalProfile(c, additionalSelects, func(client *msgraphsdkgo.GraphServiceClient, querySelects []string) (gmodels.ServicePrincipalable, error) {
		return client.ServicePrincipals().ByServicePrincipalId(ID).Get(c, &serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
				Select: querySelects,
			},
		})
	})
}

func SyncServicePrincipalProfileByAppID(c ctx.RequestContext, appID string, additionalSelects []string) (ctx.RequestContext, *ProfileDoc, gmodels.ServicePrincipalable, error) {
	return syncServicePrincipalProfile(c, additionalSelects, func(client *msgraphsdkgo.GraphServiceClient, querySelects []string) (gmodels.ServicePrincipalable, error) {
		return client.ServicePrincipalsWithAppId(&appID).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
				Select: querySelects,
			},
		})
	})
}

func syncServicePrincipalProfile(c ctx.RequestContext, additionalSelects []string,
	getSp func(client *msgraphsdkgo.GraphServiceClient, querySelects []string) (gmodels.ServicePrincipalable, error)) (ctx.RequestContext, *ProfileDoc, gmodels.ServicePrincipalable, error) {
	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	bad := func(e error) (ctx.RequestContext, *ProfileDoc, gmodels.ServicePrincipalable, error) {
		return c, nil, nil, e
	}
	if err != nil {
		return bad(err)
	}
	doc := &ProfileDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderProfile,
				NamespaceID:       NamespaceIDGraph,
				ResourceProvider:  models.ProfileResourceProviderServicePrincipal,
			},
		},
	}
	querySelects := []string{"id", "displayName", "appId", "servicePrincipalType"}
	querySelects = append(querySelects, additionalSelects...)
	sp, err := getSp(gclient, querySelects)
	if err != nil {
		err = graph.HandleMsGraphError(err)
		if errors.Is(err, graph.ErrMsGraphResourceNotFound) {
			return bad(fmt.Errorf("%w,%w", base.ErrResponseStatusNotFound, err))
		}
		return bad(err)
	}

	doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderServicePrincipal
	doc.DisplayName = sp.GetDisplayName()
	doc.ID = *sp.GetId()
	doc.AppId = sp.GetAppId()
	doc.ServicePrincipalType = sp.GetServicePrincipalType()

	_, err = resdoc.GetDocService(c).Upsert(c, doc, nil)
	return c, doc, sp, err
}
