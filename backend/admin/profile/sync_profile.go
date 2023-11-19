package profile

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
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

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
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
	switch namespaceProvider {
	case models.NamespaceProviderServicePrincipal:
		dirObj, err := gclient.ServicePrincipals().ByServicePrincipalId(namespaceId).Get(c, &serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "appId", "servicePrincipalType"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return base.ErrResponseStatusNotFound
			}
			return err
		}
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderServicePrincipal
		doc.DisplayName = dirObj.GetDisplayName()
		doc.ID = *dirObj.GetId()
		doc.AppId = dirObj.GetAppId()
		doc.ServicePrincipalType = dirObj.GetServicePrincipalType()

	case models.NamespaceProviderGroup:
		dirObj, err := gclient.Groups().ByGroupId(namespaceId).Get(c, &groups.GroupItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return base.ErrResponseStatusNotFound
			}
			return err
		}
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderGroup
		doc.DisplayName = dirObj.GetDisplayName()
		doc.ID = *dirObj.GetId()

	case models.NamespaceProviderUser:
		dirObj, err := gclient.Users().ByUserId(namespaceId).Get(c, &users.UserItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "userPrincipalName", "mail"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return base.ErrResponseStatusNotFound
			}
			return err
		}
		doc.PartitionKey.ResourceProvider = models.ProfileResourceProviderUser
		doc.DisplayName = dirObj.GetDisplayName()
		doc.ID = *dirObj.GetId()
		doc.UserPrincipalName = dirObj.GetUserPrincipalName()
		doc.Mail = dirObj.GetMail()
		// ok
	default:
		return fmt.Errorf("%w: namespace provider %s not supported", base.ErrResponseStatusBadRequest, namespaceProvider)
	}

	resp, err := resdoc.GetDocService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}
	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
