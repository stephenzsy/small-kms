package scep

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
)

type scepServer struct {
	azCredential        azcore.TokenCredential
	azKeysClient        *azkeys.Client
	smallKmsAdminClient *admin.Client

	msGraphClient                *msgraphsdk.GraphServiceClient
	msIntunesScepEndpoint        string
	msIntunesScepEndpointRefresh time.Time

	cachedCaCert   []byte
	cachedCaTime   time.Time
	fetchCaCertMtx sync.RWMutex
}

const (
	EnvVarSmallKmsAdminClientID = "SMALLKMS_ADMIN_CLIENTID"
	EnvVarSmallKmsAdminEndpoint = "SMALLKMS_ADMIN_ENDPOINT"
	SMALLKMS_API_SCOPE          = "SMALLKMS_API_SCOPE"
)

func NewScepServer() ServerInterface {
	s := scepServer{}

	var err error
	s.azCredential, err = common.GetAzCredential(os.Getenv(EnvVarSmallKmsAdminClientID))
	if err != nil {
		log.Panicf("Failed to get az credential: %s", err.Error())
	}
	apiScope := common.MustGetenv(SMALLKMS_API_SCOPE)
	s.smallKmsAdminClient, err = admin.NewClient(common.MustGetenv(EnvVarSmallKmsAdminEndpoint),
		admin.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			authHeader := req.Header.Get("Authorization")
			if len(authHeader) == 0 {
				token, tokenError := s.azCredential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{apiScope}})
				if tokenError != nil {
					return tokenError
				}
				req.Header.Set("Authorization", "Bearer "+token.Token)
			}
			return nil
		}))
	if err != nil {
		log.Panicf("Failed to create admin client: %s", err.Error())
	}

	s.azKeysClient, err = common.GetAzKeysClient()
	if err != nil {
		log.Panicf("Failed to create az keys client: %s", err.Error())
	}

	s.msGraphClient, err = msgraphsdk.NewGraphServiceClientWithCredentials(s.azCredential, []string{"Application.Read.All"})
	if err := s.refreshServiceMap(); err != nil {
		log.Panicf("Failed to refresh service map for %s", err.Error())
	}
	return &s
}
