package resdoc

import "github.com/stephenzsy/small-kms/backend/models"

func (doc *ResourceDoc) ToRef() (ref models.Ref) {
	ref.ID = doc.ID
	ref.Deleted = doc.Deleted
	ref.UpdatedBy = doc.UpdatedBy
	ref.Updated = doc.Timestamp.Time
	return
}
