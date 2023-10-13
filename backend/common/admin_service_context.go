package common

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
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

type contextKey string

const adminServerClientProviderContextKey contextKey = "adminServerClient"
const isElevatedContextKey contextKey = "isElevated"

func WithAdminServerClientProvider(c context.Context, p AdminServerClientProvider) context.Context {
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
	GetKeyvaultCertificateResourceScopeID(certificateName string, category string) string
}

const adminServerRequestClientProvierContextKey contextKey = "adminServerRequestClient"

func WithAdminServerRequestClientProvider(c RequestContext, p AdminServerRequestClientProvider) RequestContext {
	return c.WitValue(adminServerRequestClientProvierContextKey, p)
}

func GetAdminServerRequestClientProvider(c RequestContext) AdminServerRequestClientProvider {
	if p, ok := c.Value(adminServerRequestClientProvierContextKey).(AdminServerRequestClientProvider); ok {
		return p
	}
	return nil
}
