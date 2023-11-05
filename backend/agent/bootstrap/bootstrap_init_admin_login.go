package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/common"
)

type tokenCache struct {
	Account  *public.Account `json:"account,omitempty"`
	Token    []byte          `json:"tokens"`
	filename string
}

// Export implements cache.ExportReplace.
func (tc *tokenCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) (err error) {
	tc.Token, err = cache.Marshal()
	return err
}

// Replace implements cache.ExportReplace.
func (tc *tokenCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	return cache.Unmarshal(tc.Token)
}

func (tc *tokenCache) Close() {
	cacheFileBytes, _ := json.Marshal(tc)
	os.WriteFile(tc.filename, cacheFileBytes, 0640)
}

func newAppTokenCache(tokenCacheFile string) *tokenCache {
	appTokenCache := &tokenCache{
		filename: tokenCacheFile,
	}
	if tokenJson, err := os.ReadFile(tokenCacheFile); err == nil {
		json.Unmarshal(tokenJson, appTokenCache)
	}
	return appTokenCache
}

func getAppWithSharedTokenCache(c context.Context, appTokenCache *tokenCache, silent bool, forceDeviceCode bool) (*public.Client, *public.AuthResult, error) {
	bad := func(err error) (*public.Client, *public.AuthResult, error) {
		return nil, nil, err
	}

	envSvc := common.NewEnvService()
	if clientID, ok := envSvc.RequireNonWhitespace(common.EnvKeyAzClientID, common.IdentityEnvVarPrefixApp); !ok {
		return bad(envSvc.ErrMissing(common.EnvKeyAzClientID))
	} else if tenantID, ok := envSvc.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixApp); !ok {
		return bad(envSvc.ErrMissing(common.EnvKeyAzTenantID))
	} else if apiAuthScope, ok := envSvc.RequireNonWhitespace(agentcommon.EnvKeyAPIAuthScope, common.IdentityEnvVarPrefixApp); !ok {
		return bad(envSvc.ErrMissing(agentcommon.EnvKeyAPIAuthScope))
	} else {
		appClient, err := public.New(clientID,
			public.WithAuthority(fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)),
			public.WithCache(appTokenCache))
		if err != nil {
			return bad(err)
		}
		authScopes := []string{apiAuthScope}
		if appTokenCache.Account != nil {
			if authResult, err := appClient.AcquireTokenSilent(c, authScopes, public.WithTenantID(tenantID), public.WithSilentAccount(*appTokenCache.Account)); err == nil {
				return &appClient, &authResult, nil
			} else {
				fmt.Printf("Failed to acquire token silently: %v\n", err)
			}
		}
		if silent {
			return bad(errors.New("silent login failed"))
		}

		if !forceDeviceCode {
			if resp, err := appClient.AcquireTokenInteractive(c, authScopes, public.WithTenantID(tenantID),
				public.WithRedirectURI(fmt.Sprintf("msal%s://auth", clientID)),
			); err == nil {
				appTokenCache.Account = &resp.Account
				return &appClient, &resp, nil
			}
		}
		if resp, err := appClient.AcquireTokenByDeviceCode(c, authScopes, public.WithTenantID(tenantID)); err == nil {
			fmt.Printf("\033[1;33m%s\033[0m\n", resp.Result.Message)
			if r, err := resp.AuthenticationResult(c); err != nil {
				return bad(err)
			} else {
				appTokenCache.Account = &r.Account
				return &appClient, &r, nil
			}
		} else {
			return bad(err)
		}
	}
}

var _ cache.ExportReplace = (*tokenCache)(nil)

func (*ServicePrincipalBootstraper) Login(c context.Context, tokenCacheFile string, forceDeviceCode bool) error {
	if tokenCacheFile == "" {
		return errors.New("missing client cert path")
	}

	appTokenCache := newAppTokenCache(tokenCacheFile)
	defer appTokenCache.Close()

	_, _, err := getAppWithSharedTokenCache(c, appTokenCache, false, forceDeviceCode)

	return err
}
