package admin

import (
	"context"
	"encoding/json"
	"fmt"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

const (
	RefPropertyKeyDisplayName = "displayName"
	RefPropertyKeyThumbprint  = "thumbprint"
)

func PagerToList[D any](ctx context.Context, pager *azruntime.Pager[azcosmos.QueryItemsResponse]) (results []*D, err error) {
	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			err = fmt.Errorf("pager scanner error: %w", scanErr)
			return
		}
		for _, itemBytes := range t.Items {
			item := new(D)
			if err = json.Unmarshal(itemBytes, item); err != nil {
				err = fmt.Errorf("item failed to serialize: %w", scanErr)
				return
			}
			results = append(results, item)
		}
	}
	return
}

func baseDocPopulateRefWithMetadata(d kmsdoc.KmsDocument, ref *RefWithMetadata, nsType NamespaceTypeShortName) {
	if d == nil || ref == nil {
		return
	}
	ref.ID = d.GetUUID()
	ref.NamespaceID = d.GetNamespaceID()
	ref.Updated = d.GetUpdated()
	updatedById, updatedByName := d.GetUpdatedBy()
	ref.UpdatedBy = fmt.Sprintf("%s/%s", updatedById, updatedByName)
	ref.NamespaceType = nsType
	ref.Deleted = d.GetDeleted()
}

func profileDocPopulateRefWithMetadata(d graph.GraphProfileDocument, ref *RefWithMetadata, nsType NamespaceTypeShortName) {
	if d == nil || ref == nil {
		return
	}
	baseDocPopulateRefWithMetadata(d, ref, nsType)
	ref.DisplayName = d.GetDisplayName()
}
