package api

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
)

type QueryPolicyItemsParams struct {
	ExtraColumns  []string
	PolicyLocator *base.ResourceLocator
}

const (
	PolicyItemsQueryColumnCreated  = "c.iat"
	PolicyItemsQueryColumnNotAfter = "c.exp"
)

func QueryPolicyItems[DocType base.QueryDocument](c context.Context,
	partitionKey base.DocNamespacePartitionKey,
	params QueryPolicyItemsParams) *base.DocPager[DocType] {

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(params.ExtraColumns...).
		WithOrderBy(fmt.Sprintf("%s DESC", PolicyItemsQueryColumnCreated))

	if params.PolicyLocator != nil {
		qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: params.PolicyLocator.String()})
	}

	return base.NewQueryDocPager[DocType](c, qb, partitionKey)
}
