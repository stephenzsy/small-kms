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
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type configManager struct {
	envConfig         *agentcommon.AgentEnv
	configDir         RootConfigDir
	client            agentclient.ClientWithResponsesInterface
	identityProcessor *identityProcessor
}

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
		envConfig: envConfig,
		configDir: RootConfigDir{ConfigDir(configDir)},
	}
	cm.identityProcessor = &identityProcessor{cm: cm}

	certPath := cm.configDir.Config(agentmodels.AgentConfigNameIdentity).ConfigFile(configFileClientCert, true)
	if certPath == "" {
		if clientCertPath, ok := envSvc.RequireAbsPath(common.EnvKeyAzClientCertPath, common.IdentityEnvVarPrefixAgent); !ok {
			return nil, envSvc.ErrMissing(common.EnvKeyAzClientCertPath)
		} else {
			certPath = clientCertPath
		}
	}

	err = cm.configureClient(certPath)
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
