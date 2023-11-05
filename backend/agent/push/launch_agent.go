package agentpush

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func (s *agentServer) apiLaunchAgentContainer(c ctx.RequestContext, req LaunchAgentRequest) error {
	if exposedPorts, portBindings, err := nat.ParsePortSpecs(req.ExposedPortSpecs); err != nil {
		return err
	} else {
		secretMounts := make([]mount.Mount, 0, len(req.Secrets))
		for _, secret := range req.Secrets {
			secretMounts = append(secretMounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   secret.Source,
				Target:   "/run/secrets/" + secret.TargetName,
				ReadOnly: true,
			})
		}
		var networkConfig *network.NetworkingConfig
		if req.NetworkName != "" {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					req.NetworkName: {},
				},
			}
		}
		envService := s.EnvService().Clone()
		envService.SetValue("AGENT_PUSH_ENDPOINT", req.PushEndpoint)
		if req.MsEntraIdClientCertSecretName != "" {
			envService.SetValue("AZURE_CLIENT_CERTIFICATE_PATH", "/run/secrets/"+req.MsEntraIdClientCertSecretName)
		}
		result, err := s.dockerClient.ContainerCreate(c,
			&container.Config{
				ExposedPorts: exposedPorts,
				Env:          envService.Export(),
				Cmd:          []string{"/agent-server", string(req.Mode), req.ListenerAddress},
				Image:        fmt.Sprintf("%s:%s", s.acrImageRepo, req.ImageTag),
				StopSignal:   "SIGINT",
				StopTimeout:  utils.ToPtr(10),
			},
			&container.HostConfig{
				Binds:        req.HostBinds,
				PortBindings: portBindings,
				Mounts:       secretMounts,
			},
			networkConfig, nil, req.ContainerName)
		if err != nil {
			return err
		}
		if err := s.dockerClient.ContainerStart(c, result.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, result)
	}
}
