package agentconfigmanager

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/managedapp"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type echoTask struct {
	buildID      string
	configUpdate <-chan *AgentEndpointConfiguration
	newEcho      func(config *AgentEndpointConfiguration) (*echo.Echo, error)
	agentEnv     *agentcommon.AgentEnv
	endpoint     string
	mode         agentcommon.AgentSlot
}

// Name implements taskmanager.Task.
func (*echoTask) Name() string {
	return "Echo"
}

func GetTLSDefaultConfig(config *AgentEndpointConfiguration) (*tls.Config, error) {
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
	instanceIdenfier := base.IDFromUUID(instanceID)
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
				base.ID("me"), instanceIdenfier, managedapp.AgentInstanceFields{
					Version:  config.Version,
					Endpoint: et.endpoint,
					BuildID:  et.buildID,
					Mode:     managedapp.AgentMode(et.mode),
				})
		}
	}
	return nil
}

var _ taskmanager.Task = (*echoTask)(nil)

func NewEchoTask(buildID string,
	newEcho func(config *AgentEndpointConfiguration) (*echo.Echo, error),
	cm ConfigManager,
	endpoint string,
	mode agentcommon.AgentSlot) *echoTask {
	return &echoTask{
		buildID:      buildID,
		newEcho:      newEcho,
		configUpdate: cm.ConfigUpdate(),
		agentEnv:     cm.EnvConfig(),
		endpoint:     endpoint,
		mode:         mode,
	}
}