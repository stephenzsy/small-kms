package kv

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type AzKeyVaultService interface {
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
	AzSecretsClient() *azsecrets.Client
}

type internalContextKey int

const (
	AzKeyVaultServiceContextKey internalContextKey = iota
	delegatedAzSecretsClientContextKey
)

func GetAzKeyVaultService(c context.Context) AzKeyVaultService {
	if s, ok := c.Value(AzKeyVaultServiceContextKey).(AzKeyVaultService); ok {
		return s
	}
	return nil
}
