package managedapp

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile/v2"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func listManagedApps(c context.Context) ([]*ManagedAppRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder(base.ProfileResourceKindManagedApp).
		WithExtraColumns(profile.QueryColumnDisplayName, queryColumnApplicationID, queryColumnServicePrincipalID)
	storageNsID := getManageAppDocStorageNamespaceID(c)
	pager := base.NewQueryDocPager[*ManagedAppDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ManagedAppDoc) *ManagedAppRef {
		r := &ManagedAppRef{}
		d.PopulateModelRef(r)
		r.NID = storageNsID
		r.NamespaceKind = base.NamespaceKindProfile
		r.NamespaceIdentifier = namespaceIdentifierManagedApp
		r.ResourceKind = base.ProfileResourceKindManagedApp
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
