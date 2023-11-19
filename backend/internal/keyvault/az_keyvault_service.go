package kv

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stephenzsy/small-kms/backend/models"
)

type MaterialNameKind string

const (
	MaterialNameKindCertificate    MaterialNameKind = "c"
	MaterialNameKindKey            MaterialNameKind = "k"
	MaterialNameKindSecret         MaterialNameKind = "s"
	MaterialNameKindCertificateKey MaterialNameKind = "ck"
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

func GetMaterialName(
	kind MaterialNameKind,
	nsProvider models.NamespaceProvider, nsID string, policyID string) string {

	return fmt.Sprintf("%s-%s-%s-%s", kind, nsProvider, nsID, policyID)
}
