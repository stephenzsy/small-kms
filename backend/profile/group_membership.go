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
func (s *server) SyncGroupMemberOf(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, groupId uuid.UUID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceId)
	var gclient *msgraphsdkgo.GraphServiceClient
	var err error
	c, gclient, err = graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	checkMemberGroupsBody := directoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	checkMemberGroupsBody.SetGroupIds([]string{groupId.String()})
	resp, err := gclient.DirectoryObjects().ByDirectoryObjectId(string(namespaceId)).CheckMemberGroups().Post(c, checkMemberGroupsBody, nil)
	if err != nil {
		return err
	}
	memberGroupIds := resp.GetValue()
	docSvc := base.GetAzCosmosCRUDService(c)
	for _, memberGroupId := range memberGroupIds {
		memberOfDoc := new(GroupMembershipDoc)
		memberOfDoc.init(namespaceKind, namespaceId, RelTypeGroupMemberOf, base.NewNamespaceIdentifier(namespaceKind, base.IDFromUUID(uuid.MustParse(memberGroupId))))
		if err := docSvc.Upsert(c, memberOfDoc, nil); err != nil {
			return err
		}

		memberDoc := new(GroupMembershipDoc)
		memberDoc.init(base.NamespaceKindGroup, base.IDFromUUID(uuid.MustParse(memberGroupId)), RelTypeGroupMember, base.NewNamespaceIdentifier(namespaceKind, namespaceId))
		if err := docSvc.Upsert(c, memberDoc, nil); err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusNoContent)
}
