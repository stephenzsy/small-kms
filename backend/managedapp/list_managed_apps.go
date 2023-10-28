package managedapp

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type ManagedAppQueryDoc struct {
	profile.ProfileQueryDoc
	ApplicationID        uuid.UUID `json:"applicationId"`
	ServicePrincipalID   uuid.UUID `json:"servicePrincipalId"`
	ServicePrincipalType *string   `json:"servicePrincipalType,omitempty"`
}

var _ base.ModelRefPopulater[ManagedAppRef] = (*ManagedAppQueryDoc)(nil)

func (d *ManagedAppQueryDoc) PopulateModelRef(r *ManagedAppRef) {
	d.ProfileQueryDoc.PopulateModelRef(&r.ProfileRef)
	r.ApplicationID = d.ApplicationID
	r.ServicePrincipalID = d.ServicePrincipalID
}

func listManagedApps(c context.Context) ([]*ManagedAppRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(profile.QueryColumnDisplayName, queryColumnApplicationID, queryColumnServicePrincipalID)
	storageNsID := base.NewDocNamespacePartitionKey(base.NamespaceKindProfile,
		base.StringIdentifier(namespaceIDNameManagedApp),
		base.ProfileResourceKindManagedApp)
	pager := base.NewQueryDocPager[*ManagedAppQueryDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *ManagedAppQueryDoc) *ManagedAppRef {
		r := &ManagedAppRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)

}
