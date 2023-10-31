package keeper

import (
	"context"

	"github.com/rs/zerolog/log"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/base"
)

type ConfigiManagerState int

type ConfigManager struct {
	envConfig       *agentutils.AgentEnv
	configDir       string
	configProcessor AgentConfigServerProcessor
	isReady         bool
}

func (m *ConfigManager) IsReady() bool {
	return m.isReady
}

func (m *ConfigManager) LoadConfig(c context.Context) bool {
	return false
}

func (m *ConfigManager) PullConfig(c context.Context) error {
	logger := log.Ctx(c).With().Str("step", "pull config").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")
	client, err := m.envConfig.AgentClient()
	if err != nil {
		return err
	}
	resp, err := client.GetAgentConfigServerWithResponse(c, base.NamespaceKindServicePrincipal, base.StringIdentifier("me"))
	if err != nil {
		return err
	}
	_, err = m.configProcessor.ProcessUpdate(c, resp.JSON200)
	return err
}

func NewConfigManager(configDir string) (*ConfigManager, error) {
	envConfig, err := agentutils.NewAgentEnv()
	return &ConfigManager{
		envConfig: envConfig,
		configDir: configDir,
		configProcessor: AgentConfigServerProcessor{
			configDir: configDir,
			envConfig: envConfig,
		},
	}, err
}
