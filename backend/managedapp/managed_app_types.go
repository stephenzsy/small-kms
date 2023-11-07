package managedapp

import (
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
)

type ResourceReference = base.ResourceReference

type (
	managedAppRefComposed struct {
		profile.ProfileRef
		ManagedAppRefFields
	}

	agentConfigComposed struct {
		base.ResourceReference
		AgentConfigFields
	}

	agentConfigServerComposed struct {
		AgentConfig
		AgentConfigServerFields
	}

	agentInstanceComposed struct {
		base.ResourceReference
		AgentInstanceFields
	}

	agentConfigRadiusComposed struct {
		AgentConfig
		AgentConfigRadiusFields
	}
)
