package graph

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type ContextKey int

const (
	ServiceClientIDContextKey ContextKey = iota
	ServiceMsGraphClientContextKey
	ServiceMsGraphClientClientIDContextKey
	delegatedMsGraphClientContextKey
)

func GetServiceMsGraphClient(c context.Context) *msgraph.GraphServiceClient {
	if p, ok := c.Value(ServiceMsGraphClientContextKey).(*msgraph.GraphServiceClient); ok {
		return p
	}
	return nil
}

func WithDelegatedMsGraphClient(c ctx.RequestContext) (ctx.RequestContext, *msgraph.GraphServiceClient, error) {
	return auth.WithDelegatedClient[msgraph.GraphServiceClient, ContextKey](c, delegatedMsGraphClientContextKey, func(creds azcore.TokenCredential) (*msgraph.GraphServiceClient, error) {
		return msgraph.NewGraphServiceClientWithCredentials(creds, nil)
	})
}
