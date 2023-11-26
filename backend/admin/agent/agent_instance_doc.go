package agentadmin

import (
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type AgentInstanceDoc struct {
	resdoc.ResourceDoc

	Endpoint         string                         `json:"endpoint"`
	State            agentmodels.AgentInstanceState `json:"state"`
	ConfigVersion    string                         `json:"configVersion"`
	BuildID          string                         `json:"buildId"`
	TlsCertificateID string                         `json:"tlsCertificateId"`
	JwtVerfyKeyID    string                         `json:"jwtVerifyKeyId"`
}

func (doc *AgentInstanceDoc) ToRef() (m *agentmodels.AgentInstanceRef) {
	if doc == nil {
		return nil
	}
	m = &agentmodels.AgentInstanceRef{}
	m.Ref = doc.ResourceDoc.ToRef()
	m.Endpoint = doc.Endpoint
	m.State = doc.State
	m.ConfigVersion = doc.ConfigVersion
	m.BuildId = doc.BuildID
	return
}

func (doc *AgentInstanceDoc) ToModel() (m *agentmodels.AgentInstance) {
	if doc == nil {
		return nil
	}
	m = &agentmodels.AgentInstance{
		AgentInstanceRef: *doc.ToRef(),
	}
	m.TlsCertificateId = doc.TlsCertificateID
	m.JwtVerifyKeyId = doc.JwtVerfyKeyID
	return
}
