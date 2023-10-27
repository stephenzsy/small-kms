package managedapp

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listManagedApps(c context.Context) ([]*ManagedAppRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(profile.QueryColumnDisplayName, queryColumnApplicationID, queryColumnServicePrincipalID)
	storageNsID := base.NewDocNamespacePartitionKey(base.NamespaceKindProfile, namespaceIdentifierManagedApp, base.ProfileResourceKindManagedApp)
	pager := base.NewQueryDocPager[*ManagedAppDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ManagedAppDoc) *ManagedAppRef {
		r := &ManagedAppRef{}
		d.PopulateModelRef(r)

		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
