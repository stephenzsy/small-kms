package admin

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
)

type adminServer struct {
	azBlobClient          *azblob.Client
	azBlobContainerClient *azblobcontainer.Client
}

func (s *adminServer) msGraphClient(c context.Context) (*msgraph.GraphServiceClient, error) {
	if authIdentity, ok := auth.GetAuthIdentity(c); ok {
		if creds, err := authIdentity.GetOnBehalfOfTokenCredential(nil, nil); err != nil {
			return nil, err
		} else {
			return msgraph.NewGraphServiceClientWithCredentials(creds, nil)
		}
	}
	return nil, fmt.Errorf("%w: no auth header to authenticate to graph service", common.ErrStatusUnauthorized)
}

func (s *adminServer) msGraphAppClient() (*msgraph.GraphServiceClient, error) {
	return nil, nil
}
