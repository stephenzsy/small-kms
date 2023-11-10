package agentpush

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stephenzsy/small-kms/backend/base"
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
		clonedEnv := s.EnvService().Clone()
		clonedEnv.SetValue("AGENT_PUSH_ENDPOINT", req.PushEndpoint)
		if req.MsEntraIdClientCertSecretName != "" {
			clonedEnv.SetValue("AZURE_CLIENT_CERTIFICATE_PATH", "/run/secrets/"+req.MsEntraIdClientCertSecretName)
		}
		for _, reqValue := range req.Env {
			splitted := strings.SplitN(reqValue, "=", 2)
			if len(splitted) == 2 {
				clonedEnv.SetValue(splitted[0], splitted[1])
			}
		}
		result, err := s.dockerClient.ContainerCreate(c,
			&container.Config{
				ExposedPorts: exposedPorts,
				Env:          clonedEnv.Export(),
				Cmd:          []string{"/agent-server", string(req.Mode), req.ListenerAddress},
				Image:        fmt.Sprintf("%s:%s", s.acrImageRepo, req.ImageTag),
				StopSignal:   "SIGINT",
				StopTimeout:  utils.ToPtr(10),
			},
			&container.HostConfig{
				Binds:        req.HostBinds,
				PortBindings: portBindings,
				Mounts:       secretMounts,
				RestartPolicy: container.RestartPolicy{
					Name: "unless-stopped",
				},
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

func (s *agentServer) apiStopContainer(c ctx.RequestContext, containerID string) error {
	currentContainer, err := s.dockerClient.ContainerInspect(c, containerID)
	if err != nil {
		return err
	}
	if len(currentContainer.Config.Cmd) > 2 && currentContainer.Config.Cmd[1] == string(s.mode) {
		return fmt.Errorf("%w: cannot stop container of the same type: %s", base.ErrResponseStatusBadRequest, s.mode)
	}

	if err := s.dockerClient.ContainerStop(c, containerID, container.StopOptions{}); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
