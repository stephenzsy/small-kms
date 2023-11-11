package managedapp

import (
	"time"

	frconfig "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
	"github.com/stephenzsy/small-kms/backend/base"
)

type AgentConfigRadiusDoc struct {
	AgentConfigDoc
	Container AgentContainerConfiguration   `json:"container,omitempty"`
	Clients   []frconfig.RadiusClientConfig `json:"clients,omitempty"`
	DebugMode *bool                         `json:"debugMode,omitempty"`
}

func (d *AgentConfigRadiusDoc) populateModel(m *AgentConfigRadius) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.Container = &d.Container
	m.Clients = d.Clients
	m.RefreshAfter = time.Now().Add(time.Hour * 24)
	m.DebugMode = d.DebugMode
}

func (d *AgentConfigRadiusDoc) init(nsKind base.NamespaceKind, nsIdentifier base.ID) {
	d.AgentConfigDoc.init(nsKind, nsIdentifier, base.AgentConfigNameRadius)
}
