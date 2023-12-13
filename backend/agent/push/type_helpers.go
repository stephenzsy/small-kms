package agentpush

type ProxiedResponse interface {
	StatusCode() int
	GetBody() []byte
}

func (r *AgentDockerContainerInspectResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerContainerStopResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerContainerRemoveResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentLaunchAgentResponse) GetBody() []byte {
	return r.Body
}

func (r *PushAgentConfigRadiusResponse) GetBody() []byte {
	return r.Body
}
