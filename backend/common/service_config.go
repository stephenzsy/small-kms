package common

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	IdentityEnvVarNameAzTenantID          = "AZURE_TENANT_ID"
	IdentityEnvVarNameAzClientID          = "AZURE_CLIENT_ID"
	IdentityEnvVarNameAzClientSecret      = "AZURE_CLIENT_SECRET"
	IdentityEnvVarNameAzClientCertPath    = "AZURE_CLIENT_CERTIFICATE_PATH"
	IdentityEnvVarNameAzSubscriptionID    = "AZURE_SUBSCRIPTION_ID"
	IdentityEnvVarNameAzResourceGroupName = "AZURE_RESOURCE_GROUP_NAME"
	IdentityEnvVarNameUseManagedIdentity  = "USE_MANAGED_IDENTITY"

	IdentityEnvVarPrefixService = "SERVICE_"
	IdentityEnvVarPrefixApp     = "APP_"
	IdentityEnvVarPrefixAgent   = "AGENT_"
)

type AzureCredentialServiceIdentity struct {
	creds    azcore.TokenCredential
	tenantID string
}

// TokenCredential implements AzureIdentity.
func (identity AzureCredentialServiceIdentity) TokenCredential() azcore.TokenCredential {
	return identity.creds
}

func (identity AzureCredentialServiceIdentity) TenantID() string {
	return identity.tenantID
}

var _ AzureIdentity = (*AzureCredentialServiceIdentity)(nil)

func LookupEnvWithDefault(envKey string, defaultValue string) string {
	if env, ok := os.LookupEnv(envKey); ok {
		return env
	}
	return defaultValue
}

func LookupPrefixedEnvWithDefault(envVarPrefix, envKey string, defaultValue string) string {
	if env, ok := os.LookupEnv(fmt.Sprintf("%s%s", envVarPrefix, envKey)); ok {
		return env
	}
	return LookupEnvWithDefault(envKey, defaultValue)
}

func IsEnvValueTrue(envValue string) bool {
	return envValue == "true" || envValue == "1"
}

func NewAzureIdentityFromEnv(envVarPrefix string) (AzureIdentity, error) {
	if IsEnvValueTrue(LookupPrefixedEnvWithDefault(envVarPrefix, IdentityEnvVarNameUseManagedIdentity, "")) {
		opts := azidentity.ManagedIdentityCredentialOptions{}
		// use managed identity, use this option to speed up first request ttl
		clientId := LookupPrefixedEnvWithDefault(envVarPrefix, IdentityEnvVarNameAzClientID, "")
		if clientId != "" {
			opts.ID = azidentity.ClientID(clientId)
		}

		creds, err := azidentity.NewManagedIdentityCredential(&opts)
		return AzureCredentialServiceIdentity{
			creds:    creds,
			tenantID: LookupPrefixedEnvWithDefault(envVarPrefix, IdentityEnvVarNameAzTenantID, ""),
		}, err
	}
	creds, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		TenantID: LookupPrefixedEnvWithDefault(envVarPrefix, IdentityEnvVarNameAzTenantID, ""),
	})
	return AzureCredentialServiceIdentity{
		creds:    creds,
		tenantID: LookupPrefixedEnvWithDefault(envVarPrefix, IdentityEnvVarNameAzTenantID, ""),
	}, err
}
