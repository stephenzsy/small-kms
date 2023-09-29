package admin

import (
	"context"
	"encoding/json"
	"fmt"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
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

func baseDocPopulateRef(d *kmsdoc.BaseDoc, ref *Ref, nsType NamespaceTypeShortName) {
	if d == nil || ref == nil {
		return
	}
	ref.ID = d.ID.GetUUID()
	ref.NamespaceID = d.NamespaceID
	ref.Updated = d.Updated
	ref.UpdatedBy = d.UpdatedBy
	ref.NamespaceType = nsType
	ref.Deleted = d.Deleted
}
