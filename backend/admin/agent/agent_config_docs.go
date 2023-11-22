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

type AgentConfigBundleDoc struct {
	resdoc.ResourceDoc

	Items map[agentmodels.AgentConfigName]AgentConfigBundleDocItem `json:"items"`
}

type AgentConfigDoc struct {
	resdoc.ResourceDoc
	Version []byte `json:"version"`
}

type AgentConfigDocIdentity struct {
	AgentConfigDoc
	KeyCredentialsCertificatePolicyID string `json:"keyCredentialsCertificatePolicyId"`
}

func (d *AgentConfigDocIdentity) ToModel() (m agentmodels.AgentConfigIdentity) {
	m.Name = agentmodels.AgentConfigNameIdentity
	m.Updated = d.Timestamp.Time
	m.Version = hex.EncodeToString(d.Version)
	m.KeyCredentialCertificatePolicyId = d.KeyCredentialsCertificatePolicyID
	return m
}
