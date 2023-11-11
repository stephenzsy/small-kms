package radius

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/agent/configmanager"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type radiusContainerLauncher struct {
	dockerClient   dockerclient.APIClient
	acrLoginServer string
}

// After implements configmanager.ContextConfigHandler.
func (*radiusContainerLauncher) After(c context.Context) (context.Context, error) {
	return c, nil
}

// Before implements configmanager.ContextConfigHandler.
func (s *radiusContainerLauncher) Before(c context.Context) (context.Context, error) {
	config, ok := c.Value(contextKeyRadiusConfigProcessed).(*ProcessedRadiusConfig)
	if !ok || config == nil {
		return c, nil
	}
	logger := log.Ctx(c)
	exposedPorts, portBindings, err := nat.ParsePortSpecs(config.fetchedConfig.Container.ExposedPortSpecs)
	if err != nil {
		return c, err
	}
	var cmd []string
	if config.fetchedConfig.DebugMode != nil && *config.fetchedConfig.DebugMode {
		cmd = []string{"radiusd", "-X"}
	}
	if loginServer, err := acr.ExtractACRLoginServer(config.fetchedConfig.Container.ImageRepo); err != nil {
		return c, err
	} else if s.acrLoginServer != loginServer {
		return c, fmt.Errorf("image repo %s is not supported", config.fetchedConfig.Container.ImageRepo)
	}
	var networkConfig *network.NetworkingConfig
	if config.fetchedConfig.Container.NetworkName != "" {
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				config.fetchedConfig.Container.NetworkName: {},
			},
		}
	}
	// clean up previous containers
	if err := s.dockerClient.ContainerStop(c, config.fetchedConfig.Container.ContainerName, container.StopOptions{}); err != nil {
		logger.Error().Err(err).Msgf("failed to stop container %s", config.fetchedConfig.Container.ContainerName)
	}
	if err := s.dockerClient.ContainerRemove(c, config.fetchedConfig.Container.ContainerName, types.ContainerRemoveOptions{}); err != nil {
		logger.Error().Err(err).Msgf("failed to remove container %s", config.fetchedConfig.Container.ContainerName)
	}
	result, err := s.dockerClient.ContainerCreate(c,
		&container.Config{
			ExposedPorts: exposedPorts,
			Cmd:          cmd,
			Image:        fmt.Sprintf("%s:%s", config.fetchedConfig.Container.ImageRepo, config.fetchedConfig.Container.ImageTag),
			StopSignal:   "SIGINT",
			StopTimeout:  utils.ToPtr(10),
		},
		&container.HostConfig{
			Binds:        config.HostBinds,
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		networkConfig, nil, config.fetchedConfig.Container.ContainerName)
	if err != nil {
		return c, err
	}
	logger.Debug().Msgf("created container %s", result.ID)
	if err := s.dockerClient.ContainerStart(c, result.ID, types.ContainerStartOptions{}); err != nil {
		return c, err
	}
	logger.Debug().Msgf("started container %s", result.ID)
	return c, nil

}

var _ configmanager.ContextConfigHandler = (*radiusContainerLauncher)(nil)

func NewRadiusContainerLauncher(dockerClient dockerclient.APIClient, envSvc common.EnvService) (*radiusContainerLauncher, error) {
	launcher := &radiusContainerLauncher{
		dockerClient: dockerClient,
	}
	if acrImageRepo, ok := envSvc.RequireNonWhitespace(agentcommon.EnvKeyAcrImageRepository, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, envSvc.ErrMissing(agentcommon.EnvKeyAcrImageRepository)
	} else if loginServer, err := acr.ExtractACRLoginServer(acrImageRepo); err != nil {
		return nil, err
	} else {
		launcher.acrLoginServer = loginServer
	}
	return launcher, nil
}
