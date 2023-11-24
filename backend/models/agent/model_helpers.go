package agentmodels

type (
	agentConfigIdentityComposed struct {
		AgentConfigRef
		AgentConfigIdentityFields
	}

	agentConfigEndpointComposed struct {
		AgentConfigRef
		AgentConfigEndpointFields
	}
)
