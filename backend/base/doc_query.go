package base

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
)

type DocPager[D BaseDocument] struct {
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

func ToDocPager[D BaseDocument](pager *azruntime.Pager[azcosmos.QueryItemsResponse]) *DocPager[D] {
	return &DocPager[D]{innerPager: pager}
}

type CosmosQueryBuilder struct {
	Columns      []string
	WhereClauses []string
	OrderBy      string
	Parameters   []azcosmos.QueryParameter
}

func (b *CosmosQueryBuilder) BuildQuery() (string, []azcosmos.QueryParameter) {
	sb := strings.Builder{}
	sb.WriteString("SELECT ")
	for i, column := range b.Columns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(column)
	}
	sb.WriteString(" FROM c")
	for i, clause := range b.WhereClauses {
		if i == 0 {
			sb.WriteString(" WHERE (")
		} else {
			sb.WriteString(" AND (")
		}
		sb.WriteString(clause)
		sb.WriteString(")")
	}
	if b.OrderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(b.OrderBy)
	}
	return sb.String(), b.Parameters
}

func NewDefaultCosmoQueryBuilder() *CosmosQueryBuilder {
	return &CosmosQueryBuilder{
		Columns: queryDefaultColumns[:],
	}
}

func (b *CosmosQueryBuilder) WithExtraColumns(columns ...string) *CosmosQueryBuilder {
	b.Columns = append(b.Columns, columns...)
	return b
}

func (b *CosmosQueryBuilder) WithOrderBy(clause string) *CosmosQueryBuilder {
	b.OrderBy = clause
	return b
}

func NewQueryDocPager[D BaseDocument](docService AzCosmosCRUDDocService, queryBuilder *CosmosQueryBuilder, storageNamespaceID DocNamespacePartitionKey) *DocPager[D] {
	query, parameters := queryBuilder.BuildQuery()
	log.Debug().Str("query", query).Interface("parameters", parameters).Msg("NewQueryDocPager")
	pager := docService.NewQueryItemsPager(query, storageNamespaceID, &azcosmos.QueryOptions{
		QueryParameters: parameters,
	})
	return &DocPager[D]{innerPager: pager}
}
