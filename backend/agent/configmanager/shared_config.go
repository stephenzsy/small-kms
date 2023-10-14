package cm

import (
	"runtime"

	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type sharedConfig struct {
	client             agentclient.ClientWithResponsesInterface
	serviceRuntimeInfo shared.ServiceRuntimeInfo
}

func (sc *sharedConfig) AgentClient() agentclient.ClientWithResponsesInterface {
	return sc.client
}

func (sc *sharedConfig) ServiceRuntime() *shared.ServiceRuntimeInfo {
	return &sc.serviceRuntimeInfo
}

func (sc *sharedConfig) init(
	buildID string,
	client agentclient.ClientWithResponsesInterface) {
	sc.client = client
	sc.serviceRuntimeInfo = shared.ServiceRuntimeInfo{
		BuildID:   buildID,
		GoVersion: runtime.Version(),
	}
}
