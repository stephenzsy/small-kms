package profile

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func importProfile(c ctx.RequestContext, resourceKind base.ResourceKind, rID base.Identifier) (*Profile, error) {
	nsCtx := ns.GetNSContext(c)

	if !rID.IsUUID() || rID.UUID().Version() != 4 {
		return nil, fmt.Errorf("%w: resource ID of imported profile must be a valid GUID", base.ErrResponseStatusBadRequest)
	}

	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return nil, err
	}

	var doc ProfileCRUDDoc
	switch resourceKind {
	case base.ProfileResourceKindServicePrincipal:
		sp, err := gclient.ServicePrincipals().ByServicePrincipalId(rID.String()).Get(c, &serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "appId"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return nil, fmt.Errorf("%w: service principal not found: %s", base.ErrResponseStatusNotFound, rID.String())
			}
			return nil, err
		}
		d := new(ServicePrincipalProfileDoc)
		d.Init(nsCtx.Identifier(), resourceKind, rID, *sp.GetDisplayName())
		d.AppID, err = uuid.Parse(*sp.GetAppId())
		if err != nil {
			return nil, err
		}
		doc = d
	case base.ProfileResourceKindGroup:
		group, err := gclient.Groups().ByGroupId(rID.String()).Get(c, &groups.GroupItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return nil, fmt.Errorf("%w: group not found: %s", base.ErrResponseStatusNotFound, rID.String())
			}
			return nil, err
		}
		d := new(GroupProfileDoc)
		d.Init(nsCtx.Identifier(), resourceKind, rID, *group.GetDisplayName())
		if err != nil {
			return nil, err
		}
		doc = d
	case base.ProfileResourceKindUser:
		user, err := gclient.Users().ByUserId(rID.String()).Get(c, &users.UserItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "userPrincipalName"},
			},
		})
		if err != nil {
			err = base.HandleMsGraphError(err)
			if errors.Is(err, base.ErrMsGraphResourceNotFound) {
				return nil, fmt.Errorf("%w: user not found: %s", base.ErrResponseStatusNotFound, rID.String())
			}
			return nil, err
		}
		d := new(UserProfileDoc)
		d.Init(nsCtx.Identifier(), resourceKind, rID, *user.GetDisplayName())
		d.UserPrincipalName = *user.GetUserPrincipalName()
		if err != nil {
			return nil, err
		}
		doc = d
	default:
		return nil, fmt.Errorf("%w: invalid profile kind: %s", base.ErrResponseStatusBadRequest, resourceKind)
	}

	err = base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return nil, err
	}

	m := new(Profile)
	doc.PopulateModel(m)
	return m, nil
}
