package api

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stephenzsy/small-kms/backend/common"
)

type appConfidentialIdentity struct {
	tenantID               string
	clientID               string
	clientSecret           string
	clientSecretCredential *azidentity.ClientSecretCredential
}

// ClientID implements auth.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) ClientID() string {
	return i.clientID
}

// TenantID implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) TenantID() string {
	return i.tenantID
}

// GetOnBehalfOfTokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) NewOnBehalfOfTokenCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	return azidentity.NewOnBehalfOfCredentialWithSecret(i.tenantID, i.clientID, userAssertion, i.clientSecret, opts)
}

// TokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) TokenCredential() azcore.TokenCredential {
	return i.clientSecretCredential
}

var _ common.AzureAppConfidentialIdentity = (*appConfidentialIdentity)(nil)

func getAppConfidentialIdentity(envSvc common.EnvService) (*appConfidentialIdentity, error) {

	appId := appConfidentialIdentity{}
	var ok bool
	if appId.tenantID, ok = envSvc.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixApp); !ok {
		return nil, envSvc.ErrMissing(common.EnvKeyAzTenantID)
	}
	if appId.clientID, ok = envSvc.RequireNonWhitespace(common.EnvKeyAzClientID, common.IdentityEnvVarPrefixApp); !ok {
		return nil, envSvc.ErrMissing(common.EnvKeyAzClientID)
	}
	if appId.clientSecret, ok = envSvc.RequireNonWhitespace(common.EnvKeyAzClientSecret, common.IdentityEnvVarPrefixApp); !ok {
		return nil, envSvc.ErrMissing(common.EnvKeyAzClientSecret)
	}
	var err error
	if appId.clientSecretCredential, err = azidentity.NewClientSecretCredential(
		appId.tenantID, appId.clientID, appId.clientSecret, nil); err != nil {
		return nil, err
	}
	return &appId, nil
}
