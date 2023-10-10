package agentconfig

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func GetAgentConfiguration(c RequestContext, configName models.AgentConfigName) (*models.AgentConfiguration, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	switch configName {
	case models.AgentConfigNameActiveHostBootstrap:
		return handleGetAgentActiveHostBootstrap(c, nsID)
	}
	return nil, fmt.Errorf("%w: invalid step", common.ErrStatusBadRequest)
}

func PutAgentConfiguration(c RequestContext, configName models.AgentConfigName, configParams models.AgentConfigurationParameters) (*models.AgentConfiguration, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	switch configName {
	case models.AgentConfigNameActiveHostBootstrap:
		return handlePutAgentActiveHostBootstrap(c, nsID, configParams)
	}
	return nil, fmt.Errorf("%w: invalid step", common.ErrStatusBadRequest)
}
