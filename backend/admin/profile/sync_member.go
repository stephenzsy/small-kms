package profile

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// SyncMemberGroup implements admin.ServerInterface.
func (*ProfileServer) SyncMemberOf(ec echo.Context, namespaceProvider models.NamespaceProvider, nsID string, groupID string) error {
	c := ec.(ctx.RequestContext)

	nsID = ns.ResolveMeNamespace(c, nsID)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(nsID)); !authOk {
		return base.ErrResponseStatusForbidden
	}
	switch namespaceProvider {
	case models.NamespaceProviderUser:
		// ok
	default:
		return base.ErrResponseStatusBadRequest
	}

	// get user profile
	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}

	queryBody := directoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	queryBody.SetGroupIds([]string{groupID})
	resp, err := gclient.DirectoryObjects().ByDirectoryObjectId(nsID).CheckMemberGroups().Post(c, queryBody, nil)
	if err != nil {
		return err
	}
	if len(resp.GetValue()) != 1 {
		return fmt.Errorf("%w, user %s is not a member of %s", base.ErrResponseStatusBadRequest, nsID, groupID)
	}
	groupID = resp.GetValue()[0]

	// check profiles
	if _, err = syncProfileInternal(c, namespaceProvider, nsID); err != nil {
		return err
	}
	if _, err = syncProfileInternal(c, models.NamespaceProviderGroup, groupID); err != nil {
		return err
	}

	docSvc := resdoc.GetDocService(c)
	memberOfDoc := &resdoc.LinkResourceDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: namespaceProvider,
				NamespaceID:       nsID,
				ResourceProvider:  models.ResourceProviderLink,
			},
			ID: fmt.Sprintf("%s-%s", models.LinkProviderGraphMemberOf, groupID),
		},
		LinkTo:       resdoc.NewDocIdentifier(models.NamespaceProviderGroup, groupID, models.ResourceProviderLink, fmt.Sprintf("%s-%s", models.LinkProviderGraphMember, nsID)),
		LinkProvider: models.LinkProviderGraphMemberOf,
	}
	if _, err := docSvc.Upsert(c, memberOfDoc, nil); err != nil {
		return err
	}

	memberDoc := &resdoc.LinkResourceDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderGroup,
				NamespaceID:       groupID,
				ResourceProvider:  models.ResourceProviderLink,
			},
			ID: fmt.Sprintf("%s-%s", models.LinkProviderGraphMember, nsID),
		},
		LinkTo:       resdoc.NewDocIdentifier(namespaceProvider, nsID, models.ResourceProviderLink, fmt.Sprintf("%s-%s", models.LinkProviderGraphMemberOf, groupID)),
		LinkProvider: models.LinkProviderGraphMember,
	}
	if _, err := docSvc.Upsert(c, memberDoc, nil); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, memberOfDoc.ToModel())
}
