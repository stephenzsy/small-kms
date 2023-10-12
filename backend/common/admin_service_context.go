package common

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type AdminServerClientProvider interface {
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
	AzCosmosContainerClient() *azcosmos.ContainerClient
	CertsAzBlobContainerClient() *container.Client
	MsGraphClient() *msgraphsdkgo.GraphServiceClient
}

const adminServerClientProviderContextKey contextKey = "adminServerClient"

func WithAdminServerClientProvider(c RequestContext, p AdminServerClientProvider) RequestContext {
	return RequestContextWithValue(c, adminServerClientProviderContextKey, p)
}

func ContextWithAdminServerClientProvider(c context.Context, p AdminServerClientProvider) context.Context {
	return context.WithValue(c, adminServerClientProviderContextKey, p)
}

func GetAdminServerClientProvider(c context.Context) AdminServerClientProvider {
	if p, ok := c.Value(adminServerClientProviderContextKey).(AdminServerClientProvider); ok {
		return p
	}
	return nil
}

type AdminServerRequestClientProvider interface {
	MsGraphClient() (*msgraphsdkgo.GraphServiceClient, error)
	ArmRoleAssignmentsClient() (*armauthorization.RoleAssignmentsClient, error)
	GetKeyvaultCertificateResourceScopeID(certificateName string) string
}

const adminServerRequestClientProvierContextKey contextKey = "adminServerRequestClient"

func WithAdminServerRequestClientProvider(c RequestContext, p AdminServerRequestClientProvider) RequestContext {
	return RequestContextWithValue(c, adminServerRequestClientProvierContextKey, p)
}

func GetAdminServerRequestClientProvider(c RequestContext) AdminServerRequestClientProvider {
	if p, ok := c.Value(adminServerRequestClientProvierContextKey).(AdminServerRequestClientProvider); ok {
		return p
	}
	return nil
}
