package kmsdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type DocPager[D KmsDocument] struct {
	innerPager *azruntime.Pager[azcosmos.QueryItemsResponse]
}

func (p *DocPager[D]) More() bool {
	return p.innerPager.More()
}

func (p *DocPager[D]) NextPage(c context.Context) (items []D, err error) {
	t, err := p.innerPager.NextPage(c)
	if err != nil {
		return nil, err
	}
	if t.Items == nil {
		return nil, nil
	}
	items = make([]D, 0, len(t.Items))
	for i, itemBytes := range t.Items {
		if err = json.Unmarshal(itemBytes, &items[i]); err != nil {
			err = fmt.Errorf("item failed to serialize: %w", err)
			return
		}
	}
	return
}

var _ utils.ItemsPager[KmsDocument] = (*DocPager[KmsDocument])(nil)

func DefaultQueryGetWhereClause(string) string {
	return ""
}

func QueryItemsPager[D KmsDocument](
	c common.ServiceContext,
	nsID DocNsID,
	kind DocKind,
	getColumns func(baseColumns []string) []string,
	getWhereClause func(tableName string) string,
	queryParameters []azcosmos.QueryParameter) *DocPager[D] {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT ")
	columns := getColumns(getDefaultQueryColumns())
	for i, column := range columns {
		if i > 0 {
			queryBuilder.WriteString(",")
		}
		queryBuilder.WriteString("c.")
		queryBuilder.WriteString(column)
	}
	queryBuilder.WriteString(" FROM c WHERE c.namespaceId = @namespaceId AND c.kind = @kind")
	andClause := getWhereClause("c")
	if andClause != "" {
		queryBuilder.WriteString(" AND (")
		queryBuilder.WriteString(andClause)
		queryBuilder.WriteString(")")
	}

	qp := []azcosmos.QueryParameter{
		{Name: "@kind", Value: string(kind)},
		{Name: "@namespaceId", Value: nsID.String()}}
	qp = append(qp, queryParameters...)

	azPager := cc.NewQueryItemsPager(queryBuilder.String(),
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: qp,
		})

	return &DocPager[D]{innerPager: azPager}
}
