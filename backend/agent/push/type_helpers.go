package agentpush

type ProxiedResponse interface {
	StatusCode() int
	GetBody() []byte
}

func (r *AgentDockerContainerListResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerContainerInspectResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerImageListResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerInfoResponse) GetBody() []byte {
	return r.Body
}

func (r *AgentDockerNetworkListResponse) GetBody() []byte {
	return r.Body
}
