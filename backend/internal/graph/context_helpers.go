package graph

import (
	"context"

	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type ContextKey int

const (
	ServiceMsGraphClientContextKey ContextKey = iota
	delegatedMsGraphClientContextKey
)

func GetServiceMsGraphClient(c context.Context) *msgraph.GraphServiceClient {
	if p, ok := c.Value(ServiceMsGraphClientContextKey).(*msgraph.GraphServiceClient); ok {
		return p
	}
	return nil
}

func WithDelegatedMsGraphClient(c ctx.RequestContext) (ctx.RequestContext, *msgraph.GraphServiceClient, error) {
	creds, err := auth.GetAuthIdentity(c).GetOnBehalfOfTokenCredential(c, nil)
	if err != nil {
		return c, nil, err
	}
	client, err := msgraph.NewGraphServiceClientWithCredentials(creds, nil)
	if err != nil {
		return c, client, err
	}
	return c.WithValue(delegatedMsGraphClientContextKey, client), client, nil
}

func GetDelegatedMsGraphCLient(c context.Context) *msgraph.GraphServiceClient {
	if p, ok := c.Value(delegatedMsGraphClientContextKey).(*msgraph.GraphServiceClient); ok {
		return p
	}
	return nil
}
