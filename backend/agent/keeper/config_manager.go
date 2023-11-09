package keeper

import (
	"context"

	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

type ConfigiManagerState int

type ConfigManager struct {
	EnvConfig        *agentcommon.AgentEnv
	configDir        string
	configProcessor  agentConfigServerProcessor
	hasAttemptedLoad bool
}

func (m *ConfigManager) HasAttemptedLoad() bool {
	return m.hasAttemptedLoad
}

func (m *ConfigManager) LoadConfig(c context.Context) (AgentServerConfiguration, error) {
	m.hasAttemptedLoad = true
	return m.configProcessor.InitialLoad(c)
}

func (m *ConfigManager) PullConfig(c context.Context) (AgentServerConfiguration, bool, error) {
	logger := log.Ctx(c).With().Str("step", "pull config").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")
	client, err := m.EnvConfig.AgentClient()
	if err != nil {
		return nil, false, err
	}
	resp, err := client.GetAgentConfigServerWithResponse(c, base.NamespaceKindServicePrincipal, base.ID("me"))
	if err != nil || resp.StatusCode() != 200 {
		return nil, false, err
	}
	return m.configProcessor.ProcessUpdate(c, resp.JSON200)
}

func NewConfigManager(envSvc common.EnvService, mode managedapp.AgentMode) (*ConfigManager, error) {
	configDir, ok := envSvc.Require(agentcommon.EnvKeyAgentConfigDir, common.IdentityEnvVarPrefixAgent)
	if !ok {
		return nil, envSvc.ErrMissing(agentcommon.EnvKeyAgentConfigDir)
	}
	envConfig, err := agentcommon.NewAgentEnv(envSvc, mode)
	return &ConfigManager{
		EnvConfig: envConfig,
		configDir: configDir,
		configProcessor: agentConfigServerProcessor{
			configDir: configDir,
			envConfig: envConfig,
		},
	}, err
}
