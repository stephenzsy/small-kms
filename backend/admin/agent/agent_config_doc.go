package agentadmin

import (
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type AgentConfigDoc struct {
	resdoc.ResourceDoc

	KeyCrendentialsCertificatePolicyID string `json:"keyCredentialsCertificatePolicyId"`
}

func (d *AgentConfigDoc) ToModel() (m agentmodels.AgentConfig) {
	m.Ref = d.ToRef()
	m.KeyCredentialsCertificatePolicyId = d.KeyCrendentialsCertificatePolicyID
	return m
}
