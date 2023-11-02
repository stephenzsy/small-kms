package keeper

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/managedapp"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type echoTask struct {
	buildID      string
	configUpdate <-chan AgentServerConfiguration
	newEcho      func(config AgentServerConfiguration) (*echo.Echo, error)
	agentEnv     *agentutils.AgentEnv
	endpoint     string
}

// Name implements taskmanager.Task.
func (*echoTask) Name() string {
	return "Echo"
}

func GetTLSDefaultConfig(config AgentServerConfiguration) (*tls.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair(config.TLSCertificateBundleFile(), config.TLSCertificateBundleFile())
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}, nil
}

// Start implements taskmanager.Task.
func (et *echoTask) Start(c context.Context, sigCh <-chan os.Signal) error {
	logger := log.Ctx(c).With().Str("task", et.Name()).Logger()
	logger.Debug().Msg("echo server starting")
	active := true
	var e *echo.Echo
	agentClient, err := et.agentEnv.AgentClient()
	if err != nil {
		return err
	}

	instanceID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(et.endpoint))
	instanceIdenfier := base.UUIDIdentifier(instanceID)
	for active {
		select {
		case <-c.Done():
			return c.Err()
		case <-sigCh:
			active = false
			if e != nil {
				err := e.Shutdown(c)
				logger.Info().Err(err).Msg("echo server shutdown")
				return err
			}
		case config := <-et.configUpdate:
			if e != nil {
				logger.Info().Err(e.Shutdown(c)).Msg("echo server shutdown")
			}
			e, err := et.newEcho(config)
			if err != nil {
				logger.Error().Err(err).Msg("echo server failed to start")
				continue
			}
			go e.StartServer(e.TLSServer)

			agentClient.PutAgentInstance(c, base.NamespaceKindServicePrincipal,
				base.StringIdentifier("me"), instanceIdenfier, managedapp.AgentInstanceFields{
					Version:  config.GetVersion(),
					Endpoint: et.endpoint,
					BuildID:  et.buildID,
				})
		}
	}
	return nil
}

var _ taskmanager.Task = (*echoTask)(nil)

func NewEchoTask(buildID string, newEcho func(config AgentServerConfiguration) (*echo.Echo, error), keeper *keeperTaskExecutor, endpoint string) *echoTask {
	return &echoTask{
		buildID:      buildID,
		newEcho:      newEcho,
		configUpdate: keeper.ConfigUpdate(),
		agentEnv:     keeper.cm.envConfig,
		endpoint:     endpoint,
	}
}
