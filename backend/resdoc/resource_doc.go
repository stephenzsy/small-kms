package resdoc

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stephenzsy/small-kms/backend/models"
)

type ETag = azcore.ETag

// Docs are IDed by the following
// (?<partitionID>:)<namespaceProvider>:<namespaceID>:<resourceProvider>/<resourceID>
type ResourceDoc struct {
	PartitionKey PartitionKey        `json:"namespaceId"`
	ID           string              `json:"id"`
	Timestamp    *models.NumericDate `json:"_ts,omitempty"`
	ETag         *ETag               `json:"_etag,omitempty"`
	Deleted      *time.Time          `json:"deleted,omitempty"`
	UpdatedBy    string              `json:"updatedBy,omitempty"`
}

type ResourceDocument interface {
}

type ResourceQueryDocument interface {
}
