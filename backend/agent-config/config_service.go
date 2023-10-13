package agentconfig

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var configDocs map[shared.AgentConfigName]*docConfigurator[AgentConfigDocument] = map[shared.AgentConfigName]*docConfigurator[AgentConfigDocument]{
	shared.AgentConfigNameActiveHostBootstrap: newAgentActiveHostBootStrapConfigurator(),
	shared.AgentConfigNameActiveServer:        newAgentActiveServerConfigurator(),
}

func GetAgentConfiguration(c RequestContext, configName shared.AgentConfigName, params *models.GetAgentConfigurationParams, isAdmin bool) (*shared.AgentConfiguration, error) {
	nsID := ns.GetNamespaceContext(c).GetID()

	if configurator, ok := configDocs[configName]; ok {
		doc, err := configurator.readDoc(c, nsID)
		if err != nil {
			return nil, err
		}
		if params == nil {
			// admin
			return doc.toModel(true), nil
		}

		// parse token timestamp
		shouldRefresh := false
		if params.RefreshToken != nil {
			if tokenBytes, err := base64.RawURLEncoding.DecodeString(*params.RefreshToken); err == nil {
				if refreshAfter, err := time.Parse(time.RFC3339, string(tokenBytes)); err == nil {
					shouldRefresh = time.Now().After(refreshAfter)
				}
			}
		}

		if shouldRefresh {
			if c := c.Elevate(); c != nil {
				patchOps, err := configurator.eval(c, doc)
				if err != nil {
					return nil, err
				}
				if patchOps != nil {
					// can be empty
					err = kmsdoc.Patch(c, doc, *patchOps, &azcosmos.ItemOptions{
						IfMatchEtag: utils.ToPtr(doc.GetETag()),
					})
					if err != nil {
						return nil, err
					}
				}
			}
		}

		return doc.toModel(isAdmin), nil
	}

	return nil, fmt.Errorf("%w: invalid step", common.ErrStatusBadRequest)
}

func PutAgentConfiguration(c RequestContext, configName shared.AgentConfigName, configParams shared.AgentConfigurationParameters) (*shared.AgentConfiguration, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	if configurator, ok := configDocs[configName]; ok {
		if c := c.Elevate(); c != nil {
			doc, err := configurator.preparePut(c, nsID, configParams)
			if err != nil {
				return nil, err
			}
			// store
			err = kmsdoc.Upsert(c, doc)
			if err != nil {
				return nil, err
			}
			patchOps, err := configurator.eval(c, doc)
			if err != nil {
				return nil, err
			}
			if patchOps != nil {
				// can be empty
				err = kmsdoc.Patch(c, doc, *patchOps, &azcosmos.ItemOptions{
					IfMatchEtag: utils.ToPtr(doc.GetETag()),
				})
				if err != nil {
					return nil, err
				}
			}
			return doc.toModel(true), nil
		}
		// eval after put
	}
	return nil, fmt.Errorf("%w: invalid config name: %s", common.ErrStatusBadRequest, configName)
}
