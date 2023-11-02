package agentconfig

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

var configDocs map[shared.AgentConfigName]*docConfigurator[AgentConfigDocument] = map[shared.AgentConfigName]*docConfigurator[AgentConfigDocument]{
	shared.AgentConfigNameActiveHostBootstrap: newAgentActiveHostBootStrapConfigurator(),
	shared.AgentConfigNameActiveServer:        newAgentActiveServerConfigurator(),
}

func PutAgentConfiguration(c RequestContext, configName shared.AgentConfigName, configParams shared.AgentConfigurationParameters) (*shared.AgentConfiguration, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	if configurator, ok := configDocs[configName]; ok {
		c := ctx.Elevate(c)
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
	return nil, fmt.Errorf("%w: invalid config name: %s", common.ErrStatusBadRequest, configName)
}
