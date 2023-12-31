package agentadmin

import (
	"encoding/hex"
	"time"

	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type AgentConfigBundleDocItem struct {
	Updated time.Time `json:"updated"`
	Version []byte    `json:"version"`
}

func (d *AgentConfigBundleDocItem) ToModel() (m *agentmodels.AgentConfigRef) {
	if d == nil {
		return nil
	}
	m = &agentmodels.AgentConfigRef{}
	m.Updated = d.Updated
	m.Version = hex.EncodeToString(d.Version)
	return m
}

type AgentConfigBundleDoc struct {
	resdoc.ResourceDoc

	Items map[agentmodels.AgentConfigName]*AgentConfigBundleDocItem `json:"items"`
}

func (d *AgentConfigBundleDoc) ToModel() (m agentmodels.AgentConfigBundle) {
	m.Id = d.ID
	m.Expires = time.Now().Add(24 * time.Hour)
	m.Identity = d.Items[agentmodels.AgentConfigNameIdentity].ToModel()
	m.Endpoint = d.Items[agentmodels.AgentConfigNameEndpoint].ToModel()
	return m
}

type AgentConfigDoc struct {
	resdoc.ResourceDoc
	Version []byte `json:"version"`
}
