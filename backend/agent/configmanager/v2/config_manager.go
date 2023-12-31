package agentconfigmanager

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client/v2"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type ConfigManager interface {
	Client() agentclient.ClientWithResponsesInterface
	ConfigDir() RootConfigDir

	EnvConfig() *agentcommon.AgentEnv
	ConfigUpdate() <-chan *AgentEndpointConfiguration
	CryptoProvider() cryptoprovider.CryptoProvider
}

type configManager struct {
	envConfig            *agentcommon.AgentEnv
	configDir            RootConfigDir
	client               agentclient.ClientWithResponsesInterface
	identityProcessor    *identityProcessor
	endpointProcessor    *endpointProcessor
	endpointConfigUpdate chan *AgentEndpointConfiguration
	cryptoProvider       cryptoprovider.CryptoProvider
}

// CryptoProvider implements ConfigManager.
func (cm *configManager) CryptoProvider() cryptoprovider.CryptoProvider {
	return cm.cryptoProvider
}

// ConfigDir implements ConfigManager.
func (cm *configManager) ConfigDir() RootConfigDir {
	return cm.configDir
}

func (cm *configManager) Client() agentclient.ClientWithResponsesInterface {
	return cm.client
}

func (cm *configManager) EnvConfig() *agentcommon.AgentEnv {
	return cm.envConfig
}

func (cm *configManager) ConfigUpdate() <-chan *AgentEndpointConfiguration {
	return cm.endpointConfigUpdate
}

var _ ConfigManager = (*configManager)(nil)

func (cm *configManager) pullConfig(c context.Context) (expires time.Time, err error) {

	logger := log.Ctx(c)
	resp, err := cm.client.GetAgentConfigBundleWithResponse(c, "me")
	if err != nil {
		return time.Now().Add(time.Minute * 5), err
	} else if resp.StatusCode() != http.StatusOK {
		return time.Now().Add(time.Minute * 5), fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if err := cm.identityProcessor.processIdentity(c, resp.JSON200.Identity); err != nil {
		logger.Error().Err(err).Msg("failed to process identity")
	}
	if err := cm.endpointProcessor.processEndpoint(c, resp.JSON200.Endpoint); err != nil {
		logger.Error().Err(err).Msg("failed to process endpoint")
	}

	return time.Now().Add(time.Hour * 24), nil
}

func NewConfigManager(envSvc common.EnvService, slot agentcommon.AgentSlot) (*configManager, error) {
	configDir, ok := envSvc.Require(agentcommon.EnvKeyAgentConfigDir, common.IdentityEnvVarPrefixAgent)
	if !ok {
		return nil, envSvc.ErrMissing(agentcommon.EnvKeyAgentConfigDir)
	}
	if absFilepath, err := filepath.Abs(configDir); err == nil {
		configDir = absFilepath
	}
	envConfig, err := agentcommon.NewAgentEnv(envSvc, slot)
	if err != nil {
		return nil, err
	}
	cm := &configManager{
		envConfig:            envConfig,
		configDir:            RootConfigDir{ConfigDir(configDir)},
		endpointConfigUpdate: make(chan *AgentEndpointConfiguration, 1),
	}
	cm.cryptoProvider, err = cryptoprovider.NewCryptoProvider()
	if err != nil {
		return nil, err
	}
	cm.configDir.Active(agentmodels.AgentConfigNameIdentity).EnsureExist()
	cm.configDir.Active(agentmodels.AgentConfigNameEndpoint).EnsureExist()
	cm.configDir.Certs().EnsureExist()
	cm.configDir.JWKs().EnsureExist()
	cm.identityProcessor = &identityProcessor{cm: cm}
	cm.endpointProcessor = &endpointProcessor{cm: cm}
	cm.endpointProcessor.init(context.Background())

	certPath := cm.configDir.Active(agentmodels.AgentConfigNameIdentity).ConfigFile(configFileClientCert)
	if exists, err := certPath.Exists(); err != nil {
		return nil, err
	} else if !exists {
		if clientCertPath, ok := envSvc.RequireAbsPath(common.EnvKeyAzClientCertPath, common.IdentityEnvVarPrefixAgent); !ok {
			return nil, envSvc.ErrMissing(common.EnvKeyAzClientCertPath)
		} else {
			certPath = ConfigFile(clientCertPath)
		}
	}

	err = cm.configureClient(string(certPath))
	return cm, err
}

func (cm *configManager) configureClient(clientCertPath string) error {
	if apiBaseURL, ok := cm.envConfig.RequireNonWhitespace(agentcommon.EnvKeyAPIBaseURL, common.IdentityEnvVarPrefixApp); !ok {
		return fmt.Errorf("%w: %s", common.ErrMissingEnvVar, agentcommon.EnvKeyAPIBaseURL)
	} else if apiAuthScope, ok := cm.envConfig.RequireNonWhitespace(agentcommon.EnvKeyAPIAuthScope, common.IdentityEnvVarPrefixApp); !ok {
		return fmt.Errorf("%w: %s", common.ErrMissingEnvVar, agentcommon.EnvKeyAPIAuthScope)
	} else if tenantID, ok := cm.envConfig.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixAgent); !ok {
		return cm.envConfig.ErrMissing(common.EnvKeyAzTenantID)
	} else if clientID, ok := cm.envConfig.RequireNonWhitespace(common.EnvKeyAzClientID, common.IdentityEnvVarPrefixAgent); !ok {
		return cm.envConfig.ErrMissing(common.EnvKeyAzClientID)
	} else if cert, key, err := agentcommon.ParseCertificateKeyPair(clientCertPath); err != nil {
		return err
	} else if cred, err := azidentity.NewClientCertificateCredential(
		tenantID,
		clientID,
		[]*x509.Certificate{cert}, key, nil); err != nil {
		return err
	} else if cm.client, err = agentclient.NewClientWithResponses(
		apiBaseURL,
		agentclient.WithRequestEditorFn(
			agentclient.AzTokenCredentialRequestEditorFn(
				cred, policy.TokenRequestOptions{
					Scopes: []string{apiAuthScope},
				}))); err != nil {
		return err
	} else {
		cm.identityProcessor.clientCertificate = cert
	}

	return nil
}
