package profile

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"

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

	doc, err := SyncProfileInternal(c, namespaceId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.ToModel())
}

func SyncProfileInternal(c ctx.RequestContext, namespaceId string) (*ProfileDoc, error) {
	bad := func(e error) (*ProfileDoc, error) {
		return nil, e
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return bad(err)
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
		err = base.HandleMsGraphError(err)
		if errors.Is(err, base.ErrMsGraphResourceNotFound) {
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
