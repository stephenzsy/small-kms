package agentserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type ConfigLoader struct {
	identity             common.AzureIdentity
	configDir            string
	currentConfigWrapper shared.AgentConfiguration
	currentConfig        shared.AgentConfigurationAgentActiveServer
	agentClient          *agentclient.ClientWithResponses
	azKvSecretsClient    *azsecrets.Client
}

const defaultCacheSymlink = "config.json"

func newConfigLoader(
	identity common.AzureIdentity,
	apiEndpoint string,
	apiScope string,
	tenantID string,
	configDir string,
) (ConfigLoader, error) {
	l := ConfigLoader{
		identity:  identity,
		configDir: configDir,
	}
	client, err := agentclient.NewClientWithCreds(apiEndpoint, identity.TokenCredential(), []string{apiScope}, tenantID)
	if err != nil {
		return l, err
	}
	l.agentClient = client
	if l.loadFromFile() == nil {
		l.pullCertificates(context.Background())
	}
	return l, nil
}

func (cl *ConfigLoader) pullCertificates(c context.Context) error {

	// pull certificates
	cert, err := cl.agentClient.GetCertificateWithResponse(c, shared.NamespaceKindServicePrincipal, shared.StringIdentifier("me"),
		*cl.currentConfig.ServerCertificateId, nil)
	if err != nil {
		return err
	}
	log.Info().RawJSON("cert", cert.Body).Msg("cert")
	return nil
}

func (cl *ConfigLoader) loadFromFile() error {
	configPath := filepath.Join(cl.configDir, defaultCacheSymlink)
	if contentBytes, err := os.ReadFile(configPath); err != nil {
		return err
	} else {
		err = json.Unmarshal(contentBytes, &cl.currentConfigWrapper)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cl *ConfigLoader) refreshConfig(c context.Context) (bool, error) {
	resp, err := cl.agentClient.GetAgentConfigurationWithResponse(c, shared.NamespaceKindServicePrincipal, shared.StringIdentifier("me"), shared.AgentConfigNameActiveServer, &agentclient.GetAgentConfigurationParams{
		RefreshToken: cl.currentConfigWrapper.NextRefreshToken,
		//XSmallkmsIfVersionNotMatch: &cl.currentConfigWrapper.Version,
	})
	if err != nil {
		return false, err
	}

	if resp.StatusCode() == http.StatusOK {
		// new config
		if cl.currentConfigWrapper.Version != resp.JSON200.Version {
			cl.currentConfigWrapper = *resp.JSON200
			cl.currentConfig, err = cl.currentConfigWrapper.Config.AsAgentConfigurationAgentActiveServer()
		}
		if err != nil {
			return false, err
		}
	} else if resp.StatusCode() == http.StatusNoContent {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	nextConfigPath := filepath.Join(cl.configDir, "config."+cl.currentConfigWrapper.Version+".json")
	if err := os.WriteFile(nextConfigPath, resp.Body, 0644); err != nil {
		return true, err
	}
	err = os.WriteFile(filepath.Join(cl.configDir, defaultCacheSymlink), resp.Body, 0644)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (cl *ConfigLoader) Start(c context.Context, reloadCh chan<- string) {
	for {
		var nextDuration time.Duration
		if cl.currentConfigWrapper.NextRefreshAfter == nil ||
			time.Now().After(*cl.currentConfigWrapper.NextRefreshAfter) {
			log.Debug().Msgf("refreshing config")
			refreshed, err := cl.refreshConfig(c)
			if err != nil {
				log.Error().Err(err).Msg("failed to refresh config, retry in 5 minutes")
				nextDuration = time.Duration(5 * time.Minute)
			} else {
				log.Debug().Msgf("refreshed: %v", refreshed)
				cl.pullCertificates(context.Background())
				reloadCh <- cl.currentConfigWrapper.Version
			}
		}
		if nextDuration == 0 {
			if cl.currentConfigWrapper.NextRefreshAfter != nil {
				log.Info().Msgf("next refresh after: %s", cl.currentConfigWrapper.NextRefreshAfter.String())
				nextDuration = time.Until(*cl.currentConfigWrapper.NextRefreshAfter) + time.Duration(101*time.Second)
			} else {
				nextDuration = time.Duration(5 * time.Minute)
			}
		}

		select {
		case <-c.Done():
			return
		case <-time.After(nextDuration):
			// continue the loop
			continue
		}
	}
}
