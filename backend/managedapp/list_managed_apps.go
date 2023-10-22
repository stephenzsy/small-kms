package managedapp

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listManagedApps(c context.Context) ([]*ManagedApp, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder(base.ResourceKindManagedApp).
		WithExtraColumns(queryColumnDisplayName, queryColumnApplicationID, queryColumnServicePrincipalID)
	pager := base.NewQueryDocPager[*ManagedAppDoc](docService, qb, getManageAppDocStorageNamespaceID(c))

	modelPager := utils.NewMappedItemsPager(pager, managedAppDocToModel)
	return utils.PagerToSlice(c, modelPager)

}
