package cert

import (
	"slices"

	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
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

func (s *server) allowAdminDelegatedServicePrincipal(policyNsID base.ID) authz.AuthZFunc {
	return func(c ctx.RequestContext) (ctx.RequestContext, authz.AuthzResult) {
		logger := log.Ctx(c)
		identity := auth.GetAuthIdentity(c)
		if !identity.HasAdminRole() {
			return c, authz.AuthzResultNone
		}
		c, gclient, err := graph.WithDelegatedMsGraphClient(c)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get delegated graph client")
			return c, authz.AuthzResultNone
		}
		requestAppID := identity.AppID()
		if requestAppID == "" {
			return c, authz.AuthzResultNone
		}
		sp, err := gclient.ServicePrincipalsWithAppId(&requestAppID).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id"},
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to get service principal by app id")
			return c, authz.AuthzResultNone
		}
		if *sp.GetId() == string(policyNsID) {
			return c, authz.AuthzResultAllow
		}
		return c, authz.AuthzResultNone
	}
}

func (s *server) allowGraphGroupMemberOf(nsKind base.NamespaceKind, policyNsID base.ID) authz.AuthZFunc {
	if nsKind != base.NamespaceKindGroup {
		return nil
	}
	return func(c ctx.RequestContext) (ctx.RequestContext, authz.AuthzResult) {
		logger := log.Ctx(c)
		identity := auth.GetAuthIdentity(c)
		c, gclient, err := graph.WithDelegatedMsGraphClient(c)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get delegated graph client")
			return c, authz.AuthzResultNone
		}
		dirObjBuilder := gclient.DirectoryObjects().ByDirectoryObjectId(identity.ClientPrincipalID().String())
		dirObj, err := dirObjBuilder.Get(c, &directoryobjects.DirectoryObjectItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &directoryobjects.DirectoryObjectItemRequestBuilderGetQueryParameters{
				Select: []string{"id"},
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to get directory object")
			return c, authz.AuthzResultNone
		}

		requestBody := directoryobjects.NewItemCheckMemberGroupsPostRequestBody()
		requestBody.SetGroupIds([]string{string(policyNsID)})
		resp, err := dirObjBuilder.CheckMemberGroups().Post(c, requestBody, nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to check member groups")
			return c, authz.AuthzResultNone
		}
		if !slices.Contains(resp.GetValue(), string(policyNsID)) {
			return c, authz.AuthzResultNone
		}
		switch dirObj.(type) {
		case gmodels.Userable:
			c = ns.WithNSContext(c, base.NamespaceKindUser, base.IDFromUUID(identity.ClientPrincipalID()))
		case gmodels.ServicePrincipalable:
			c = ns.WithNSContext(c, base.NamespaceKindServicePrincipal, base.IDFromUUID(identity.ClientPrincipalID()))
		default:
			logger.Error().Msgf("unsupported requester type: %s", *dirObj.GetOdataType())
			return c, authz.AuthzResultNone
		}
		return c, authz.AuthzResultAllow
	}
}
