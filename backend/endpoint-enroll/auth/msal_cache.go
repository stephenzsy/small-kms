package auth

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stephenzsy/small-kms/backend/common"
)

type cachedTokenCredential struct {
	Token    string    `json:"token"`
	Expires  time.Time `json:"exp"`
	loaded   bool
	hasToken bool
	next     azcore.TokenCredential
	filename string
}

type CachedTokenCredential interface {
	azcore.TokenCredential
	Clear() error
}

// GetToken implements exported.TokenCredential.
func (c *cachedTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !c.loaded {
		if b, err := os.ReadFile(c.filename); err == nil {
			if err = json.Unmarshal(b, c); err == nil {
				c.hasToken = true
			}
		}
		c.loaded = true
	}
	if c.hasToken && c.Expires.After(time.Now()) {
		return azcore.AccessToken{Token: c.Token, ExpiresOn: c.Expires}, nil
	}
	token, err := c.next.GetToken(ctx, options)
	if err != nil {
		return token, err
	}
	c.Token = token.Token
	c.Expires = token.ExpiresOn
	if b, err := json.Marshal(c); err == nil {
		os.WriteFile(c.filename, b, 0600)
	}
	return token, nil
}

func (c *cachedTokenCredential) Clear() error {
	return os.Remove(c.filename)
}

func newCachedTokenCredentialFromFile(filename string, next azcore.TokenCredential) *cachedTokenCredential {
	return &cachedTokenCredential{
		filename: filename,
		next:     next,
	}
}

// should only be used for install
func GetCachedTokenCredential(next azcore.TokenCredential) CachedTokenCredential {
	filename := common.GetEnvWithDefault("TOKEN_CACHE_FILE", "token-cache.json")
	return newCachedTokenCredentialFromFile(filename, next)
}
