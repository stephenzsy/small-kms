package common

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	IdentityEnvVarPrefixService = "SERVICE_"
	IdentityEnvVarPrefixApp     = "APP_"
	IdentityEnvVarPrefixAgent   = "AGENT_"
)

type AzureCredentialServiceIdentity struct {
	creds    azcore.TokenCredential
	tenantID string
	clientID string
}

// TokenCredential implements AzureIdentity.
func (identity AzureCredentialServiceIdentity) TokenCredential() azcore.TokenCredential {
	return identity.creds
}

func (identity AzureCredentialServiceIdentity) TenantID() string {
	return identity.tenantID
}

func (identity AzureCredentialServiceIdentity) ClientID() string {
	return identity.clientID
}

var _ AzureIdentity = (*AzureCredentialServiceIdentity)(nil)

func LookupEnvWithDefault(envKey string, defaultValue string) string {
	if env, ok := os.LookupEnv(envKey); ok {
		return env
	}
	return defaultValue
}

func IsEnvValueTrue(envValue string) bool {
	return envValue == "true" || envValue == "1"
}

func NewAzureIdentityFromEnv(envService EnvService, envVarPrefix string) (AzureIdentity, error) {
	if IsEnvValueTrue(envService.Default(envKeyUseManagedIdentity, "", envVarPrefix)) {
		opts := azidentity.ManagedIdentityCredentialOptions{}
		// use managed identity, use this option to speed up first request ttl
		clientId := envService.Default(EnvKeyAzClientID, "", envVarPrefix)
		if clientId != "" {
			opts.ID = azidentity.ClientID(clientId)
		}

		creds, err := azidentity.NewManagedIdentityCredential(&opts)
		return AzureCredentialServiceIdentity{
			creds:    creds,
			tenantID: envService.Default(EnvKeyAzTenantID, "", envVarPrefix),
			clientID: envService.Default(EnvKeyAzClientID, "", envVarPrefix),
		}, err
	}
	creds, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		TenantID: envService.Default(EnvKeyAzTenantID, "", envVarPrefix),
	})
	return AzureCredentialServiceIdentity{
		creds:    creds,
		tenantID: envService.Default(EnvKeyAzTenantID, "", envVarPrefix),
		clientID: envService.Default(EnvKeyAzClientID, "", envVarPrefix),
	}, err
}
