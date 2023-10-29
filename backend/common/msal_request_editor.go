package common

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func ToAzTokenCredentialRequestEditorFn(cred azcore.TokenCredential, tokenOptions policy.TokenRequestOptions) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if req.Header.Get("Authorization") == "" {
			t, err := cred.GetToken(ctx, tokenOptions)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Bearer "+t.Token)
		}
		return nil
	}
}

func ToSilenTokenRequestEditorFn(pubClient *public.Client, tokenScope string, account public.Account) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if req.Header.Get("Authorization") == "" {
			authResult, err := pubClient.AcquireTokenSilent(ctx, []string{tokenScope}, public.WithSilentAccount(account))
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Bearer "+authResult.AccessToken)
		}
		return nil
	}
}
