package common

import (
	ctx "context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type AzureIdentity interface {
	TokenCredential() azcore.TokenCredential
}

type AzureAppConfidentialIdentity interface {
	AzureIdentity
	GetOnBehalfOfTokenCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error)
}

type ServiceIdentityProvider interface {
	ServiceIdentity() AzureIdentity
}

type ConfidentialAppIdentityProvider interface {
	ConfidentialAppIdentity() AzureAppConfidentialIdentity
}

type ServerContext interface {
	ctx.Context
	ServiceIdentity() AzureIdentity
}
