package agentpush

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
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
	s.dockerClient.ContainerCreate(c,
		&container.Config{
			ExposedPorts: params.ExposedPorts,
			Env:          fromContainer.Config.Env,
			Cmd:          []string{"/agent-server", "server", params.ListenerAddress},
			Image:        params.Image,
			StopSignal:   "SIGINT",
			StopTimeout:  fromContainer.Config.StopTimeout,
		},
		&container.HostConfig{
			Binds:        fromContainer.HostConfig.Binds,
			PortBindings: params.PortBindings,
			AutoRemove:   true,
			Mounts:       fromContainer.HostConfig.Mounts,
		},
		&network.NetworkingConfig{},
		nil, fromContainer.Name+"-sidecar")

	return nil
}
