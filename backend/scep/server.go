package scep

import (
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"golang.org/x/time/rate"

	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
)

type scepRateLimits struct {
	msGraphServiceMapping *rate.Limiter

	getCa *rate.Limiter
}

type scepServer struct {
	adminServer admin.AdminServerInternal

	azCredential azcore.TokenCredential
	azKeysClient *azkeys.Client

	msGraphClient         *msgraphsdk.GraphServiceClient
	msIntunesScepEndpoint string

	rateLimiters    scepRateLimits
	getCaCertCaches map[uuid.UUID]*getCaCertCache
}

const (
	EnvVarMsGraphClientID = "MSGRAPH_CLIENT_ID"
)

var intranetNamespaceID = uuid.MustParse(string(admin.WellKnownNamespaceIDStrIntCASCEPIntranet))

func NewScepServer() ServerInterface {
	s := scepServer{
		adminServer: admin.NewAdminServer(),
		rateLimiters: scepRateLimits{
			msGraphServiceMapping: rate.NewLimiter(rate.Every(time.Hour*2), 1),
			getCa:                 rate.NewLimiter(rate.Every(time.Minute), 2),
		},
		getCaCertCaches: make(map[uuid.UUID]*getCaCertCache),
	}

	var err error
	s.azCredential, err = common.GetAzCredential(os.Getenv(EnvVarMsGraphClientID))
	if err != nil {
		log.Panicf("Failed to get az credential: %v", err)
	}

	s.azKeysClient, err = common.GetAzKeysClient()
	if err != nil {
		log.Panicf("Failed to create az keys client: %v", err)
	}

	// reserve one token for initialization
	s.msGraphClient, err = msgraphsdk.NewGraphServiceClientWithCredentials(s.azCredential, []string{"https://graph.microsoft.com/.default"})
	s.rateLimiters.msGraphServiceMapping.Allow()
	if err := s.refreshServiceMap(); err != nil {
		log.Panicf("Failed to refresh service map: %v", err)
	}
	return &s
}
