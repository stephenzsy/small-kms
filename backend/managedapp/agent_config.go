package managedapp

import (
	"github.com/stephenzsy/small-kms/backend/base"
)

type AgentConfigDoc struct {
	base.BaseDoc

	Version string `json:"version"`
}

func (d *AgentConfigDoc) init(nsKind base.NamespaceKind, nsIdentifier base.ID, configName base.NamespaceConfigName) {
	d.BaseDoc.Init(nsKind, nsIdentifier, base.ResourceKindNamespaceConfig, base.ID(configName))
}

func (d *AgentConfigDoc) PopulateModel(m *AgentConfig) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.Version = d.Version
}
