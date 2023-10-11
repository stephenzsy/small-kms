package acr

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

type authenticationClient struct {
	endpoint string
	pl       runtime.Pipeline
	aadCreds azcore.TokenCredential
	tenantID string
}

type AuthenticationClientOptions struct {
	azcore.ClientOptions
	TenantID string
}

func NewAuthenticationClient(endpoint string, aadCreds azcore.TokenCredential, options *AuthenticationClientOptions) *authenticationClient {
	if options == nil {
		options = &AuthenticationClientOptions{}
	}

	pipeline := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &options.ClientOptions)

	client := &authenticationClient{
		endpoint: endpoint,
		pl:       pipeline,
		aadCreds: aadCreds,
		tenantID: options.TenantID,
	}
	return client
}

const (
	moduleName      = "azcontainerregistry"
	moduleVersion   = "v0.2.1"
	defaultAudience = "https://containerregistry.azure.net"
	aadAudience     = "https://management.core.windows.net/"
)

func (client *authenticationClient) exchangeAADAccessTokenForACRRefreshTokenCreateRequest(ctx context.Context, service, accessToken string) (*policy.Request, error) {
	urlPath := "/oauth2/exchange"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	formData := url.Values{}

	if err != nil {
		return nil, err
	}
	formData.Set("grant_type", "access_token")
	formData.Set("service", service)
	formData.Set("tenant", client.tenantID)
	/*
		if options != nil && options.Tenant != nil {
			formData.Set("tenant", *options.Tenant)
		}*/
	formData.Set("access_token", accessToken)
	body := streaming.NopCloser(strings.NewReader(formData.Encode()))
	return req, req.SetBody(body, "application/x-www-form-urlencoded")
}

type acrAccessToken struct {
	// The access token for performing authenticated requests
	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

func (client *authenticationClient) ExchagneAADTokenForACRRefreshToken(ctx context.Context, service string) (acrAccessToken, error) {
	accessToken, err := client.aadCreds.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{aadAudience + ".default"},
	})
	if err != nil {
		return acrAccessToken{}, err
	}

	// // first find chanllenge
	// chanllengeReq, err := client.getChallengeRequest(ctx, accessToken.Token)
	// if err != nil {
	// 	return acrAccessToken{}, err
	// }
	// chanllengResp, err := client.pl.Do(chanllengeReq)
	// var service, scope string
	// if chanllengResp.StatusCode == http.StatusUnauthorized {
	// 	if service, scope, err = findServiceAndScope(chanllengResp); err != nil {
	// 		return acrAccessToken{}, err
	// 	}
	// } else {
	// 	return acrAccessToken{}, err
	// }

	req, err := client.exchangeAADAccessTokenForACRRefreshTokenCreateRequest(ctx, service, accessToken.Token)
	if err != nil {
		return acrAccessToken{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return acrAccessToken{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return acrAccessToken{}, runtime.NewResponseError(resp)
	}
	return client.exchangeAADAccessTokenForACRRefreshTokenHandleResponse(resp)
}

func (client *authenticationClient) exchangeAADAccessTokenForACRRefreshTokenHandleResponse(resp *http.Response) (acrAccessToken, error) {
	result := acrAccessToken{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return acrAccessToken{}, err
	}
	return result, nil
}

// func (client *authenticationClient) getChallengeRequest(ctx context.Context, accessToken string) (*policy.Request, error) {
// 	urlPath := "/v2/"
// 	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
// 	if err != nil {
// 		return nil, err
// 	}
// 	//req.Raw().Header["Authorization"] = []string{"Bearer " + accessToken}
// 	return req, err
// }

// func findServiceAndScope(resp *http.Response) (string, string, error) {
// 	authHeader := resp.Header.Get("WWW-Authenticate")
// 	if authHeader == "" {
// 		return "", "", errors.New("response has no WWW-Authenticate header for challenge authentication")
// 	}
// 	log.Println("Www-Authenticate: ", authHeader)

// 	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
// 	parts := strings.Split(authHeader, "\",")
// 	valuesMap := map[string]string{}
// 	for _, part := range parts {
// 		subParts := strings.Split(part, "=")
// 		if len(subParts) == 2 {
// 			valuesMap[subParts[0]] = strings.ReplaceAll(subParts[1], "\"", "")
// 		}
// 	}

// 	if _, ok := valuesMap["service"]; !ok {
// 		return "", "", errors.New("could not find a valid service in the WWW-Authenticate header")
// 	}

// 	if _, ok := valuesMap["scope"]; !ok {
// 		return "", "", errors.New("could not find a valid scope in the WWW-Authenticate header")
// 	}

// 	return valuesMap["service"], valuesMap["scope"], nil
// }
