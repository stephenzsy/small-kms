package managedapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
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

func apiGetSystemApp(c ctx.RequestContext, systemAppName SystemAppName) error {
	doc, err := apiGetSystemAppDoc(c, systemAppName)
	if err != nil {
		return err
	}
	m := &ManagedApp{}
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}

func apiGetSystemAppDoc(c ctx.RequestContext, systemAppName SystemAppName) (*ManagedAppDoc, error) {
	appID, err := resolveSystemAppID(c, systemAppName)
	if err != nil {
		return nil, err
	}
	doc := &ManagedAppDoc{}
	docService := base.GetAzCosmosCRUDService(c)
	err = docService.Read(c, base.NewDocLocator(base.NamespaceKindProfile,
		base.IDFromString(namespaceIDNameSystemApp),
		base.ProfileResourceKindManagedApp,
		base.IDFromUUID(appID)), doc, nil)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: system app not found: %s", base.ErrResponseStatusNotFound, systemAppName)
		}
		return nil, err
	}
	return doc, nil
}
