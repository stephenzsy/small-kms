package cert

import (
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
)

func (s *server) allowGeneralNonAdminAuth(c ctx.RequestContext,
	nsKind base.NamespaceKind, namespaceIdentifier base.ID) (ctx.RequestContext, ns.NSContext, error) {

	c, nsCtx := ns.WithResovingMeNSContext(c, nsKind, namespaceIdentifier)
	c, authOk := authz.Authorize(c, authz.AllowAdmin, nsCtx.AllowSelf(), s.allowGroupMember(nsKind, nsCtx.ID()))
	if !authOk {
		return c, nsCtx, base.ErrResponseStatusForbidden
	}
	return c, nsCtx, nil
}

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
			if err := docSvc.Read(c, base.NewDocLocator(base.NamespaceKindGroup, groupID,
				base.ResourceKindNamespaceConfig, base.IDFromUUID(identity.ClientPrincipalID())), relDoc, nil); err != nil || relDoc.RelType != profile.RelTypeGroupMember {
				// no membership
				return c, authz.AuthzResultNone
			}
			c = c.WithValue(groupMemberOfContextKey, groupUUID)
			return c, authz.AuthzResultAllow
		}
	}
}

// func (s *server) allowAdminDelegatedServicePrincipal(policyNsID base.ID) authz.AuthZFunc {
// 	return func(c ctx.RequestContext) (ctx.RequestContext, authz.AuthzResult) {
// 		logger := log.Ctx(c)
// 		identity := auth.GetAuthIdentity(c)
// 		if !identity.HasAdminRole() {
// 			return c, authz.AuthzResultNone
// 		}
// 		c, gclient, err := graph.WithDelegatedMsGraphClient(c)
// 		if err != nil {
// 			logger.Error().Err(err).Msg("failed to get delegated graph client")
// 			return c, authz.AuthzResultNone
// 		}
// 		requestAppID := identity.AppID()
// 		if requestAppID == "" {
// 			return c, authz.AuthzResultNone
// 		}
// 		sp, err := gclient.ServicePrincipalsWithAppId(&requestAppID).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
// 			QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
// 				Select: []string{"id"},
// 			},
// 		})
// 		if err != nil {
// 			logger.Error().Err(err).Msg("failed to get service principal by app id")
// 			return c, authz.AuthzResultNone
// 		}
// 		if *sp.GetId() == string(policyNsID) {
// 			return c.WithValue(selfGraphObjectContextKey, sp), authz.AuthzResultAllow
// 		}
// 		return c, authz.AuthzResultNone
// 	}
//}

//
