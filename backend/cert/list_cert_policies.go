package cert

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListCertPolicies implements ServerInterface.
func (s *server) ListCertPolicies(ec echo.Context, nsKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	c, nsCtx := ns.WithResovingMeNSContext(c, nsKind, namespaceIdentifier)
	c, authOk := authz.Authorize(c, authz.AllowAdmin, nsCtx.AllowSelf(), s.allowGroupMember(nsKind, nsCtx.ID()))
	if !authOk {
		return base.ErrResponseStatusForbidden
	}

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy)
	pager := base.NewQueryDocPager[*CertPolicyQueryDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertPolicyQueryDoc) *CertPolicyRef {
		r := &CertPolicyRef{}
		d.PopulateModelRef(r)
		return r
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

type CertPolicyQueryDoc struct {
	base.QueryBaseDoc
	DisplayName string `json:"displayName"`
}

// PopulateModelRef implements base.ModelRefPopulater.
func (d *CertPolicyQueryDoc) PopulateModelRef(r *CertPolicyRef) {
	if d == nil || r == nil {
		return
	}
	d.QueryBaseDoc.PopulateModelRef(&r.ResourceReference)
	r.DisplayName = d.DisplayName
}

var _ base.ModelRefPopulater[CertPolicyRef] = (*CertPolicyQueryDoc)(nil)

type internalContextKey int

const (
	groupMemberOfContextKey internalContextKey = iota
)

func (s *server) allowGroupMember(nsKind base.NamespaceKind, groupID base.ID) authz.AuthZFunc {
	if nsKind != base.NamespaceKindGroup {
		return nil
	}
	if groupUUID, ok := groupID.AsUUID(); !ok {
		return nil
	} else {
		return func(c ctx.RequestContext) (ctx.RequestContext, authz.AuthzResult) {
			identity := auth.GetAuthIdentity(c)
			docSvc := base.GetAzCosmosCRUDService(c)
			relDoc := new(profile.GroupMembershipDoc)
			if err := docSvc.Read(c, base.NewDocFullIdentifier(base.NamespaceKindGroup, groupID,
				base.ResourceKindNamespaceConfig, base.IDFromUUID(identity.ClientPrincipalID())), relDoc, nil); err != nil || relDoc.RelType != profile.RelTypeGroupMember {
				// no membership
				return c, authz.AuthzResultNone
			}
			c = c.WithValue(groupMemberOfContextKey, groupUUID)
			return c, authz.AuthzResultAllow
		}
	}
}
