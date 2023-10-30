package keeper

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/base"
)

type ConfigiManagerState int

type ConfigManager struct {
	envConfig agentutils.AgentEnv
	configDir string
	isReady   bool
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
	agentConfigServer := resp.JSON200
	logger.Debug().Any("config", agentConfigServer).Msg("agent config")
	_ = agentConfigServer

	// pull certificate
	{
		resp, err := client.GetCertificateWithResponse(c, base.NamespaceKindServicePrincipal, base.StringIdentifier("me"), agentConfigServer.TlsCertificateId)
		if err != nil {
			return err
		}
		cert := resp.JSON200
		logger.Debug().Any("cert", cert).Msg("cert")

		azSecretsClient, err := m.envConfig.AzSecretsClient()
		if err != nil {
			return err
		}
		sid := azsecrets.ID(*cert.KeyVaultSecretID)
		getSecretResposne, err := azSecretsClient.GetSecret(c, sid.Name(), sid.Version(), nil)
		if err != nil {
			return err
		}
		logger.Debug().Any("secrets resp", getSecretResposne).Msg("secrets resp")
	}

	return nil
}

func NewConfigManager(configDir string) (*ConfigManager, error) {
	envConfig, err := agentutils.NewAgentEnv()
	return &ConfigManager{
		envConfig: envConfig,
		configDir: configDir,
	}, err
}
