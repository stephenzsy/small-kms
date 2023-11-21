package profile

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
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
