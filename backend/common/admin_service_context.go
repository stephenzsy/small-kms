package common

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type AdminServerClientProvider interface {
	AzCosmosContainerClient() *azcosmos.ContainerClient
	CertsAzBlobContainerClient() *container.Client
	MsGraphClient() *msgraphsdkgo.GraphServiceClient
}

type contextKey string

const (
	AdminServerClientProviderContextKey       contextKey = "adminServerClient"
	AdminServerRequestClientProvierContextKey contextKey = "adminServerRequestClient"
)

func GetAdminServerClientProvider(c context.Context) AdminServerClientProvider {
	return c.Value(AdminServerClientProviderContextKey).(AdminServerClientProvider)
}

type AdminServerRequestClientProvider interface {
	MsGraphClient() (*msgraphsdkgo.GraphServiceClient, error)
	ArmRoleAssignmentsClient() (*armauthorization.RoleAssignmentsClient, error)
	GetKeyvaultCertificateResourceScopeID(certificateName string, category string) string
}

func GetAdminServerRequestClientProvider(c context.Context) AdminServerRequestClientProvider {
	return c.Value(AdminServerRequestClientProvierContextKey).(AdminServerRequestClientProvider)
}
