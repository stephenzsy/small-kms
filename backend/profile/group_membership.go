package profile

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

const (
	RelTypeGroupMemberOf = "groupMemberOf"
	RelTypeGroupMember   = "groupMember"
)

type GroupMembershipDoc struct {
	base.BaseDoc
	RelType string                   `json:"relType"`
	Target  base.NamespaceIdentifier `json:"target"`
}

func (d *GroupMembershipDoc) init(nsKind base.NamespaceKind, nsID base.ID, relType string, target base.NamespaceIdentifier) {
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindNamespaceConfig, target.ID())
	d.RelType = relType
	d.Target = target
}

// SyncGroupMemberOf implements ServerInterface.
func (s *server) SyncGroupMemberOf(ec echo.Context, namespaceKind base.NamespaceKind, namespaceID base.ID, groupId uuid.UUID) error {
	c := ec.(ctx.RequestContext)

	nsUUID := auth.ResolveSelfNamespace(c, string(namespaceID))
	if !auth.AuthorizeSelfOrAdmin(c, nsUUID) {
		return s.RespondRequireAdmin(c)
	}
	nsID := base.IDFromUUID(nsUUID)

	c = ns.WithNSContext(c, namespaceKind, nsID)
	var gclient *msgraphsdkgo.GraphServiceClient
	var err error
	c, gclient, err = graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	checkMemberGroupsBody := directoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	checkMemberGroupsBody.SetGroupIds([]string{groupId.String()})
	resp, err := gclient.DirectoryObjects().ByDirectoryObjectId(string(nsID)).CheckMemberGroups().Post(c, checkMemberGroupsBody, nil)
	if err != nil {
		return err
	}
	memberGroupIds := resp.GetValue()
	docSvc := base.GetAzCosmosCRUDService(c)
	for _, memberGroupId := range memberGroupIds {
		memberOfDoc := new(GroupMembershipDoc)
		memberOfDoc.init(namespaceKind, nsID, RelTypeGroupMemberOf, base.NewNamespaceIdentifier(namespaceKind, base.IDFromUUID(uuid.MustParse(memberGroupId))))
		if err := docSvc.Upsert(c, memberOfDoc, nil); err != nil {
			return err
		}

		memberDoc := new(GroupMembershipDoc)
		memberDoc.init(base.NamespaceKindGroup, base.IDFromUUID(uuid.MustParse(memberGroupId)), RelTypeGroupMember, base.NewNamespaceIdentifier(namespaceKind, nsID))
		if err := docSvc.Upsert(c, memberDoc, nil); err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// ListGroupMemberOf implements ServerInterface.
func (s *server) ListGroupMemberOf(ec echo.Context, nsKind base.NamespaceKind, namespaceID base.ID) error {
	c := ec.(ctx.RequestContext)
	nsUUID := auth.ResolveSelfNamespace(c, string(namespaceID))
	if !auth.AuthorizeSelfOrAdmin(c, nsUUID) {
		return s.RespondRequireAdmin(c)
	}
	nsID := base.IDFromUUID(nsUUID)
	c = ns.WithNSContext(c, nsKind, nsID)
	qb := base.NewDefaultCosmoQueryBuilder().WithExtraColumns(QueryColumnDisplayName).
		WithWhereClauses("c.relType = '" + RelTypeGroupMemberOf + "'")
	pager := base.NewQueryDocPager[*GroupMembershipDoc](c, qb, base.NewDocNamespacePartitionKey(nsKind, nsID, base.ResourceKindNamespaceConfig))

	modelsPager := utils.NewSerializableItemsPager(
		utils.NewMappedItemsPager(pager, func(item *GroupMembershipDoc) *base.ResourceReference {
			model := new(base.ResourceReference)
			item.PopulateModelRef(model)
			return model
		}))
	return c.JSON(http.StatusOK, modelsPager)
}
