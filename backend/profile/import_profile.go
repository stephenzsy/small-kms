package profile

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
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
