package managedapp

import (
	"time"

	"github.com/stephenzsy/small-kms/backend/agent/radius"
	"github.com/stephenzsy/small-kms/backend/base"
)

type AgentConfigRadiusDoc struct {
	AgentConfigDoc
	GlobalRadiusServerACRImageRef string                      `json:"acrImageRef,omitempty"`
	Clients                       []radius.RadiusClientConfig `json:"clients,omitempty"`
}

func (d *AgentConfigRadiusDoc) populateModel(m *AgentConfigRadius) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.AzureACRImageRef = &d.GlobalRadiusServerACRImageRef
	m.RefreshAfter = time.Now().Add(time.Hour * 24)
	m.Clients = d.Clients
}

func (d *AgentConfigRadiusDoc) init(nsKind base.NamespaceKind, nsIdentifier base.ID) {
	d.AgentConfigDoc.init(nsKind, nsIdentifier, base.AgentConfigNameRadius)
}
