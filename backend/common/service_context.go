package common

import (
	ctx "context"

	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type AzureIdentity = auth.AzureIdentity
type AzureAppConfidentialIdentity = auth.AzureAppConfidentialIdentity

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
