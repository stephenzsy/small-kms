package agentmodels

import "github.com/stephenzsy/small-kms/backend/models"

type (
	agentConfigIdentityComposed struct {
		AgentConfigRef
		AgentConfigIdentityFields
	}

	agentConfigEndpointComposed struct {
		AgentConfigRef
		AgentConfigEndpointFields
	}

	agentInstanceRefComposed struct {
		models.Ref
		AgentInstanceRefFields
	}
)
