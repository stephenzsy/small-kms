package managedapp

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

func getManagedApp(c context.Context, appID uuid.UUID) (*ManagedAppDoc, error) {

	doc := &ManagedAppDoc{}
	docService := base.GetAzCosmosCRUDService(c)
	err := docService.Read(c, base.NewDocLocator(base.NamespaceKindProfile,
		base.IDFromString(namespaceIDNameManagedApp),
		base.ProfileResourceKindManagedApp,
		base.IDFromUUID(appID)), doc, nil)
	return doc, err
}
