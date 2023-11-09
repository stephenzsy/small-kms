package managedapp

import (
	"time"

	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
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

// NextPullAfter implements configmanager.VersionedConfig.
func (cfg agentConfigComposed) NextPullAfter() time.Time {
	return cfg.RefreshAfter
}

// GetVersion implements configmanager.VersionedConfig.
func (cfg AgentConfig) GetVersion() string {
	return cfg.Version
}

var _ configmanager.VersionedConfig = (*AgentConfig)(nil)
