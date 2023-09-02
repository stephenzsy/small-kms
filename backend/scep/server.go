package scep

import (
	"log"
	"time"

	"github.com/google/uuid"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"golang.org/x/time/rate"

	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/scep/msintune"
)

type scepRateLimits struct {
	msGraphServiceMapping *rate.Limiter

	getCa *rate.Limiter
}

type scepServer struct {
	common.CommonConfig
	adminServer admin.AdminServerInternal

	msGraphClient        *msgraphsdk.GraphServiceClient
	msIntuneScepEndpoint string
	msIntuneClient       *msintune.Client

	rateLimiters    scepRateLimits
	getCaCertCaches map[uuid.UUID]*getCaCertCache
}

const (
	EnvVarMsGraphClientID = "MSGRAPH_CLIENT_ID"
)

var intranetNamespaceID = uuid.MustParse(string(admin.WellKnownNamespaceIDStrIntCASCEPIntranet))

func NewScepServer() ServerInterface {
	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic(err)
	}
	s := scepServer{
		CommonConfig: commonConfig,
		adminServer:  admin.NewAdminServer(),
		rateLimiters: scepRateLimits{
			msGraphServiceMapping: rate.NewLimiter(rate.Every(time.Hour*2), 1),
			getCa:                 rate.NewLimiter(rate.Every(time.Minute), 2),
		},
		getCaCertCaches: make(map[uuid.UUID]*getCaCertCache),
	}

	// reserve one token for initialization
	s.msGraphClient, err = msgraphsdk.NewGraphServiceClientWithCredentials(commonConfig.DefaultAzCredential(), nil)
	if err != nil {
		log.Panicf("Failed to create ms graph client: %v", err)
	}
	s.rateLimiters.msGraphServiceMapping.Allow()
	if s.msIntuneClient, err = s.refreshServiceMap(); err != nil {
		log.Panicf("Failed to refresh service map: %v", err)
	}
	return &s
}
