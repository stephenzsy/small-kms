package kv

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func WithDelegatedAzSecretsClient(c ctx.RequestContext, keyvaultEndpoint string) (ctx.RequestContext, *azsecrets.Client, error) {
	return auth.WithDelegatedClient[azsecrets.Client, internalContextKey](
		c, delegatedAzSecretsClientContextKey, func(creds azcore.TokenCredential) (*azsecrets.Client, error) {
			return azsecrets.NewClient(keyvaultEndpoint, creds, nil)
		})
}
