package admin

import (
	"context"
	"encoding/json"
	"fmt"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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
