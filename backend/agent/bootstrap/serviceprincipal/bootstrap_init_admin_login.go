package serviceprincipal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/stephenzsy/small-kms/backend/common"
)

type tokenCache struct {
	LocalAccountID string            `json:"localAccountId"`
	Tokens         map[string][]byte `json:"tokens"`
}

// Export implements cache.ExportReplace.
func (tc *tokenCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) (err error) {
	tc.Tokens[hints.PartitionKey], err = cache.Marshal()
	return err
}

// Replace implements cache.ExportReplace.
func (tc *tokenCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	if b, ok := tc.Tokens[hints.PartitionKey]; ok {
		return cache.Unmarshal(b)
	}
	return nil
}

var _ cache.ExportReplace = (*tokenCache)(nil)

func (*ServicePrincipalBootstraper) Login(c context.Context, tokenCacheFile string, forceDeviceCode bool) error {
	if tokenCacheFile == "" {
		return errors.New("missing client cert path")
	}

	appTokenCache := &tokenCache{
		Tokens: map[string][]byte{},
	}
	if tokenJson, err := os.ReadFile(tokenCacheFile); err != nil {
		json.Unmarshal(tokenJson, appTokenCache)
	}
	defer func() {
		cacheFileBytes, _ := json.Marshal(appTokenCache)
		os.WriteFile(tokenCacheFile, cacheFileBytes, 0400)
	}()

	var appClient public.Client
	var err error

	if clientID := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientID, ""); clientID == "" {
		return errors.New("missing APP_AZURE_CLIENT_ID")
	} else if tenantID := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzTenantID, ""); tenantID == "" {
		return errors.New("missing APP_AZURE_TENANT_ID")
	} else if apiAuthScope := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "APP_API_AUTH_SCOPE", ""); apiAuthScope == "" {
		return errors.New("missing APP_API_AUTH_SCOPE")
	} else {
		appClient, err = public.New(clientID,
			public.WithAuthority(fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)),
			public.WithCache(appTokenCache))
		if err != nil {
			return err
		}
		authScopes := []string{apiAuthScope}
		if appTokenCache.LocalAccountID != "" {
			if _, err := appClient.AcquireTokenSilent(c, authScopes, public.WithTenantID(tenantID), public.WithSilentAccount(public.Account{
				LocalAccountID: appTokenCache.LocalAccountID,
			})); err == nil {
				return nil
			}
		}
		if !forceDeviceCode {
			if resp, err := appClient.AcquireTokenInteractive(c, authScopes, public.WithTenantID(tenantID)); err == nil {
				appTokenCache.LocalAccountID = resp.Account.LocalAccountID
				return nil
			}
		}
		if resp, err := appClient.AcquireTokenByDeviceCode(c, authScopes); err != nil {
			fmt.Printf("\033[1;33m%s\033[0m\n", resp.Result.Message)
			if r, err := resp.AuthenticationResult(c); err != nil {
				appTokenCache.LocalAccountID = r.Account.LocalAccountID
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
