package agentconfigmanager

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/taskmanager"
)

type echoTask struct {
	buildID       string
	newEcho       func(config *AgentEndpointConfiguration) (*echo.Echo, error)
	configManager ConfigManager
	endpoint      string
	mode          agentcommon.AgentSlot
}

// Name implements taskmanager.Task.
func (*echoTask) Name() string {
	return "Echo"
}

func getActiveTLSCertificateBundleFile(cm ConfigManager) string {
	return string(cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert))
}

func GetTLSDefaultConfig(configManager ConfigManager) (*tls.Config, error) {
	certBundle := getActiveTLSCertificateBundleFile(configManager)
	tlsCert, err := tls.LoadX509KeyPair(certBundle, certBundle)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"h2", "http/1.1"},
	}, nil
}

// Start implements taskmanager.Task.
func (et *echoTask) Start(c context.Context, sigCh <-chan os.Signal) error {
	logger := log.Ctx(c).With().Str("task", et.Name()).Logger()
	logger.Debug().Msg("echo server starting")
	active := true
	var e *echo.Echo

	var lastConfig *AgentEndpointConfiguration
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
			if lastConfig != nil {
				if resp, err := et.configManager.Client().UpdateAgentInstanceWithResponse(c, "me", agentmodels.AgentInstanceParameters{
					BuildId:       et.buildID,
					Endpoint:      et.endpoint,
					ConfigVersion: lastConfig.Version,
					State:         agentmodels.AgentInstanceStateStopped,
				}); err != nil {
					logger.Error().Err(err).Msg("failed to update agent instance")
				} else if resp.StatusCode() >= 400 {
					logger.Error().Int("status", resp.StatusCode()).Msg("failed to update agent instance")
				}
			}
		case lastConfig = <-et.configManager.ConfigUpdate():
			if e != nil {
				logger.Info().Err(e.Shutdown(c)).Msg("echo server shutdown")
			}
			e, err := et.newEcho(lastConfig)
			if err != nil {
				logger.Error().Err(err).Msg("echo server failed to start")
				continue
			}
			go e.StartServer(e.TLSServer)

			if resp, err := et.configManager.Client().UpdateAgentInstanceWithResponse(c, "me", agentmodels.AgentInstanceParameters{
				BuildId:          et.buildID,
				Endpoint:         et.endpoint,
				ConfigVersion:    lastConfig.Version,
				State:            agentmodels.AgentInstanceStateRunning,
				TlsCertificateId: lastConfig.TLSCertificateID,
				JwtVerifyKeyId:   lastConfig.VerifyJwkID,
			}); err != nil {
				logger.Error().Err(err).Msg("failed to update agent instance")
			} else if resp.StatusCode() >= 400 {
				logger.Error().Int("status", resp.StatusCode()).RawJSON("body", resp.Body).Msg("failed to update agent instance")
			}
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
		buildID:       buildID,
		newEcho:       newEcho,
		configManager: cm,
		endpoint:      endpoint,
		mode:          mode,
	}
}
