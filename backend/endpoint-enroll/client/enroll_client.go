package client

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func NewClientWithCreds(server string, creds azcore.TokenCredential, scopes []string, tenantID string) (*Client, error) {
	return NewClient(server, WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		token, err := creds.GetToken(ctx, policy.TokenRequestOptions{
			Scopes:   scopes,
			TenantID: tenantID,
		})
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token.Token)
		return nil
	}))
}
