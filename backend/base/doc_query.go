package base

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type DocPager[D QueryDocument] struct {
	innerPager *azruntime.Pager[azcosmos.QueryItemsResponse]
	queryCtx   context.Context
}

func (p *DocPager[D]) More() bool {
	return p.innerPager.More()
}

func (p *DocPager[D]) NextPage() (items []D, err error) {
	t, err := p.innerPager.NextPage(p.queryCtx)
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

var _ utils.ItemsPager[QueryDocument] = (*DocPager[QueryDocument])(nil)

func ToDocPager[D QueryDocument](pager *azruntime.Pager[azcosmos.QueryItemsResponse]) *DocPager[D] {
	return &DocPager[D]{innerPager: pager}
}

type CosmosQueryBuilder struct {
	Columns           []string
	WhereClauses      []string
	OrderBy           string
	Parameters        []azcosmos.QueryParameter
	OffsetLimitClause string
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
	if b.OffsetLimitClause != "" {
		sb.WriteString(b.OffsetLimitClause)
	}
	return sb.String(), b.Parameters
}

var queryDefaultColumns = []string{
	"c.id",
	"c._ts",
	"c.deleted",
}

type QueryBaseDoc struct {
	ID        Identifier      `json:"id"`
	Timestamp jwt.NumericDate `json:"_ts"`
	Deleted   *time.Time      `json:"deleted"`
}

// PopulateModelRef implements ModelRefPopulater.
func (d *QueryBaseDoc) PopulateModelRef(r *ResourceReference) {
	if d == nil || r == nil {
		return
	}
	r.Id = d.ID
	r.Deleted = d.Deleted
	r.Updated = d.Timestamp.Time
}

// GetID implements QueryDocument.
func (d *QueryBaseDoc) GetID() identifier {
	return d.ID
}

type QueryDocument interface {
	GetID() Identifier
}

var _ QueryDocument = (*QueryBaseDoc)(nil)

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

func (b *CosmosQueryBuilder) WithOffsetLimit(offset uint, limit uint) *CosmosQueryBuilder {
	b.OffsetLimitClause = fmt.Sprintf(" OFFSET %d LIMIT %d", offset, limit)
	return b
}

func NewQueryDocPager[D QueryDocument](c context.Context, queryBuilder *CosmosQueryBuilder, storageNamespaceID DocNamespacePartitionKey) *DocPager[D] {
	query, parameters := queryBuilder.BuildQuery()
	log.Ctx(c).Debug().Str("query", query).Interface("parameters", parameters).Msg("NewQueryDocPager")
	pager := GetAzCosmosCRUDService(c).NewQueryItemsPager(query, storageNamespaceID, &azcosmos.QueryOptions{
		QueryParameters: parameters,
	})
	return &DocPager[D]{innerPager: pager,
		queryCtx: c}
}

var _ ModelRefPopulater[ResourceReference] = (*QueryBaseDoc)(nil)
