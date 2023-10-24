package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

type GraphProfileDoc struct {
	ProfileDoc
}

// GetStorageID implements base.CRUDDocHasCustomStorageID.
func (d *GraphProfileDoc) GetStorageID(context.Context) uuid.UUID {
	return d.ResourceIdentifier.UUID()
}

var _ base.CRUDDocHasCustomStorageID = (*GraphProfileDoc)(nil)

type ServicePrincipalProfileDoc struct {
	GraphProfileDoc

	AppID uuid.UUID `json:"appId"`
}
