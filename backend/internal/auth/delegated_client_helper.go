package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func WithDelegatedClient[C, CK any](c ctx.RequestContext, contextKey CK, getClient func(azcore.TokenCredential) (*C, error)) (ctx.RequestContext, *C, error) {
	if p, ok := c.Value(contextKey).(*C); ok {
		return c, p, nil
	}
	creds, err := GetAuthIdentity(c).GetOnBehalfOfTokenCredential(c, nil)
	if err != nil {
		return c, nil, err
	}
	client, err := getClient(creds)
	return c.WithValue(contextKey, client), client, err
}

func GetDelegateClient[C, CK any](c context.Context, contextKey CK) *C {
	if p, ok := c.Value(contextKey).(*C); ok {
		return p
	}
	return nil
}
