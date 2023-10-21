package graph

import (
	"context"

	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type ContextKey int

const (
	ServiceMsGraphClientContextKey ContextKey = iota
)

func GetServiceMsGraphClient(c context.Context) *msgraph.GraphServiceClient {
	if p, ok := c.Value(ServiceMsGraphClientContextKey).(*msgraph.GraphServiceClient); ok {
		return p
	}
	return nil
}
