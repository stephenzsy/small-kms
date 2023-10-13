package agentserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type ConfigLoader struct {
	baseUrl       string
	authScope     string
	identity      common.AzureIdentity
	tenantID      string
	cacheFileName string

	currentConfigWrapper shared.AgentConfiguration
	currentConfig        shared.AgentConfigurationAgentActiveServer
}

func (cl *ConfigLoader) refreshConfig(c context.Context) (bool, error) {
	log.Debug().Msgf("Base url: %s", cl.baseUrl)
	client, err := agentclient.NewClientWithCreds(cl.baseUrl, cl.identity.TokenCredential(), []string{cl.authScope}, cl.tenantID)
	if err != nil {
		return false, err
	}
	resp, err := client.AgentGetConfigurationWithResponse(c, shared.AgentConfigNameActiveServer, &agentclient.AgentGetConfigurationParams{
		RefreshToken: cl.currentConfigWrapper.NextRefreshToken,
		//XSmallkmsIfVersionNotMatch: &cl.currentConfigWrapper.Version,
	})
	if err != nil {
		return false, err
	}

	if resp.StatusCode() == http.StatusOK {
		// new config
		cl.currentConfigWrapper = *resp.JSON200
		cl.currentConfig, err = cl.currentConfigWrapper.Config.AsAgentConfigurationAgentActiveServer()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if resp.StatusCode() == http.StatusNoContent {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
}

func (cl *ConfigLoader) Start(c context.Context, reloadCh chan<- string) {
	if contentBytes, err := os.ReadFile(cl.cacheFileName); err != nil {
		log.Error().Err(err).Msgf("failed to read cache file: %s", cl.cacheFileName)
	} else {
		err = json.Unmarshal(contentBytes, &cl.currentConfigWrapper)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unmarshal cache file: %s", cl.cacheFileName)
		}
		reloadCh <- cl.currentConfigWrapper.Version
	}
	for {
		nextDuration := time.Duration(1 * time.Minute)
		if cl.currentConfigWrapper.NextRefreshAfter == nil ||
			time.Now().After(*cl.currentConfigWrapper.NextRefreshAfter) {
			log.Debug().Msgf("refreshing config")
			refreshed, err := cl.refreshConfig(c)
			if err != nil {
				log.Error().Err(err).Msg("failed to refresh config, retry in 5 minutes")
			}
			log.Debug().Msgf("refreshed: %v", refreshed)
			if refreshed {
				reloadCh <- cl.currentConfigWrapper.Version
				nextDuration = time.Until(*cl.currentConfigWrapper.NextRefreshAfter) + time.Duration(2*time.Minute)
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
