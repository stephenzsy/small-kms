package agentclient

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func AzTokenCredentialRequestEditorFn(
	certCred azcore.TokenCredential,
	tokenOptions policy.TokenRequestOptions) RequestEditorFn {

	return func(ctx context.Context, req *http.Request) error {
		if req.Header.Get("Authorization") == "" {
			t, err := certCred.GetToken(ctx, tokenOptions)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Bearer "+t.Token)
		}
		return nil
	}
}
