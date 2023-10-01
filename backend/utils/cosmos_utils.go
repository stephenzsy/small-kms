package utils

import (
	"context"
	"encoding/json"
	"fmt"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func PagerToList[D any](ctx context.Context, pager *azruntime.Pager[azcosmos.QueryItemsResponse]) (results []D, err error) {

	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			err = fmt.Errorf("pager scanner error: %w", scanErr)
			return
		}
		slice := make([]D, len(t.Items))
		for i, itemBytes := range t.Items {
			if err = json.Unmarshal(itemBytes, &slice[i]); err != nil {
				err = fmt.Errorf("item failed to serialize: %w", scanErr)
				return
			}
		}
		results = append(results, slice...)
	}
	return
}
