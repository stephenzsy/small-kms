package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func putProfile(c context.Context, resourceKind base.ResourceKind, rID base.ID, params *ProfileParameters) (*Profile, error) {
	nsCtx := ns.GetNSContext(c)

	doc := new(ProfileDoc)

	displayName := string(rID)
	if params.DisplayName != nil && *params.DisplayName != "" {
		displayName = *params.DisplayName
	}

	doc.Init(nsCtx.ID(), resourceKind, rID, displayName)
	err := base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return nil, err
	}

	m := new(Profile)
	doc.PopulateModel(m)
	return m, nil
}
