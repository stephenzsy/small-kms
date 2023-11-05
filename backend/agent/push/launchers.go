package agentpush

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type Parameters struct {
	ExposedPorts    nat.PortSet
	ListenerAddress string
	Image           string
	PortBindings    nat.PortMap
}

func (s *agentServer) LaunchSidecar(c context.Context, fromContainerID string,
	params Parameters) error {
	fromContainer, err := s.dockerClient.ContainerInspect(c, fromContainerID)
	if err != nil {
		return err
	}
	os.Environ()
	s.dockerClient.ContainerCreate(c,
		&container.Config{
			ExposedPorts: params.ExposedPorts,
			Env:          fromContainer.Config.Env,
			Cmd:          []string{"/agent-server", "server", params.ListenerAddress},
			Image:        params.Image,
			StopSignal:   "SIGINT",
			StopTimeout:  utils.ToPtr(10),
		},
		&container.HostConfig{
			Binds:        fromContainer.HostConfig.Binds,
			PortBindings: params.PortBindings,
			AutoRemove:   true,
			Mounts:       fromContainer.HostConfig.Mounts,
		},
		&network.NetworkingConfig{
			EndpointsConfig: fromContainer.NetworkSettings.Networks,
		},
		nil, fromContainer.Name+"-sidecar")

	return nil
}
