package profile

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listProfiles(c context.Context, namespaceIdentifier base.Identifier, resourceKind base.ResourceKind) ([]*ProfileRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder(resourceKind).
		WithExtraColumns(QueryColumnDisplayName)
	storageNsID := getProfileDocStorageNamespaceID(c, namespaceIdentifier)
	pager := base.NewQueryDocPager[*ProfileDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ProfileDoc) *ProfileRef {
		r := &ProfileRef{}
		d.PopulateModelRef(r)
		r.NID = storageNsID
		r.NamespaceKind = base.NamespaceKindProfile
		r.NamespaceIdentifier = namespaceIdentifier
		r.ResourceKind = resourceKind
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
