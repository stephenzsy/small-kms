package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func putProfile(c context.Context, resourceKind base.ResourceKind, rID base.Identifier, params *ProfileParameters) (*Profile, error) {
	nsCtx := ns.GetNSContext(c)

	doc := new(ProfileDoc)

	displayName := rID.String()
	if params.DisplayName != nil && *params.DisplayName != "" {
		displayName = *params.DisplayName
	}

	doc.Init(nsCtx.Identifier(), resourceKind, rID, displayName)
	err := base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return nil, err
	}

	m := new(Profile)
	doc.PopulateModel(m)
	return m, nil
}
