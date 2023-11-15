package systemapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	appadmin "github.com/stephenzsy/small-kms/backend/admin/app"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func (s *SystemAppAdminServer) GetSystemApp(ec echo.Context, id string) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return echo.ErrForbidden
	}

	appName, err := validateSystemAppName(id)
	if err != nil {
		return err
	}

	doc, appID, err := GetSystemAppDoc(c, appName)
	if err != nil {
		if errors.Is(err, echo.ErrNotFound) {
			return c.JSON(200, &models.ApplicationByAppId{
				Ref: models.Ref{
					ID: appID.String(),
				},
			})
		}
	}

	return c.JSON(200, doc.ToApplicationByAppId())
}

func resolveSystemAppID(c context.Context, systemAppName SystemAppName) (uuid.UUID, error) {
	switch systemAppName {
	case SystemAppNameBackend:
		if systemAppID, ok := c.Value(graph.ServiceClientIDContextKey).(string); ok {
			return uuid.Parse(systemAppID)
		}
	case SystemAppNameAPI:
		if systemAppID, ok := c.Value(graph.ServiceMsGraphClientClientIDContextKey).(string); ok {
			return uuid.Parse(systemAppID)
		}
	}
	return uuid.Nil, fmt.Errorf("%w: system app not found: %s", echo.ErrNotFound, systemAppName)
}

func GetSystemAppDoc(c ctx.RequestContext, systemAppName SystemAppName) (*SystemAppDoc, uuid.UUID, error) {
	appID, err := resolveSystemAppID(c, systemAppName)
	if err != nil {
		return nil, appID, err
	}
	doc := &SystemAppDoc{}
	err = resdoc.GetDocService(c).Read(c, resdoc.DocIdentifier{
		PartitionKey: resdoc.PartitionKey{
			NamespaceProvider: models.NamespaceProviderProfile,
			NamespaceID:       appadmin.AppNamespaceID,
			ResourceProvider:  models.ProfileResourceProviderSystem,
		},
		ID: appID.String(),
	}, doc, nil)
	if err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return nil, appID, fmt.Errorf("%w: system app not found: %s", echo.ErrNotFound, systemAppName)
		}
		return nil, appID, err
	}
	return doc, appID, nil
}
