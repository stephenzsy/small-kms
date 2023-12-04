package resdoc

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type ETag = azcore.ETag

// Docs are IDed by the following
// (?<partitionID>:)<namespaceProvider>:<namespaceID>:<resourceProvider>/<resourceID>
type ResourceDoc struct {
	PartitionKey PartitionKey     `json:"namespaceId"`
	ID           string           `json:"id"`
	Timestamp    *jwt.NumericDate `json:"_ts,omitempty"`
	ETag         *ETag            `json:"_etag,omitempty"`
	Deleted      *time.Time       `json:"deleted,omitempty"`
	UpdatedBy    string           `json:"updatedBy,omitempty"`
}

// GetETag implements ResourceDocument.
func (doc *ResourceDoc) GetETag() *azcore.ETag {
	return doc.ETag
}

func (doc *ResourceDoc) Identifier() DocIdentifier {
	return DocIdentifier{
		PartitionKey: doc.PartitionKey,
		ID:           doc.ID,
	}
}

// setTimestamp implements ResourceDocument.
func (doc *ResourceDoc) setTimestamp(t time.Time) {
	doc.Timestamp = jwt.NewNumericDate(t)
}

// partitionKey implements ResourceDocument.
func (doc *ResourceDoc) partitionKey() PartitionKey {
	return doc.PartitionKey
}

// setETag implements ResourceDocument.
func (d *ResourceDoc) setETag(etag azcore.ETag) {
	if d.ETag == nil {
		return
	}
	d.ETag = &etag
}

func (d *ResourceDoc) getUpdatedBy() string {
	return d.UpdatedBy
}

func (d *ResourceDoc) setUpdatedBy(value string) {
	d.UpdatedBy = value
}

func (d *ResourceDoc) prepareForWrite(c context.Context) {
	d.UpdatedBy = auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	// clear read-only fields
	d.ETag = nil
	d.Timestamp = nil
}

func (d *ResourceDoc) getID() string {
	return d.ID
}

type ResourceDocument interface {
	Identifier() DocIdentifier
	getID() string
	partitionKey() PartitionKey
	setETag(etag azcore.ETag)
	GetETag() *azcore.ETag
	prepareForWrite(c context.Context)
	setTimestamp(t time.Time)
	getUpdatedBy() string
	setUpdatedBy(value string)
}

var _ ResourceDocument = (*ResourceDoc)(nil)

type ResourceQueryDocument interface {
}
