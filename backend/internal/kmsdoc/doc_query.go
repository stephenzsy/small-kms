package kmsdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/models"
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
	items = make([]D, len(t.Items))
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

type CosmosQueryBuilder struct {
	ExtraColumns      []string
	ExtraWhereClauses []string
	OrderBy           string
	ExtraParameters   []azcosmos.QueryParameter
}

func (b *CosmosQueryBuilder) BuildQuery(kind models.ResourceKind) (string, []azcosmos.QueryParameter) {
	sb := strings.Builder{}
	sb.WriteString("SELECT ")
	for i, column := range queryDefaultColumns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("c.")
		sb.WriteString(column)
	}
	for _, column := range b.ExtraColumns {
		sb.WriteString(",c.")
		sb.WriteString(column)
	}
	sb.WriteString(" FROM c WHERE c.kind = @kind")
	for _, clause := range b.ExtraWhereClauses {
		sb.WriteString(" AND (")
		sb.WriteString(clause)
		sb.WriteString(")")
	}
	if b.OrderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(b.OrderBy)
	}
	return sb.String(), append([]azcosmos.QueryParameter{
		{Name: "@kind", Value: string(kind)}}, b.ExtraParameters...)
}

func QueryItemsPager[D KmsDocument](
	c RequestContext,
	nsID docNsIDType,
	kind models.ResourceKind,
	getQueryBuilder func(tableName string) CosmosQueryBuilder) *DocPager[D] {
	cc := c.ServiceClientProvider().AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	qb := getQueryBuilder("c")
	query, queryParameters := qb.BuildQuery(kind)

	azPager := cc.NewQueryItemsPager(query,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: queryParameters,
		})

	return &DocPager[D]{innerPager: azPager}
}
