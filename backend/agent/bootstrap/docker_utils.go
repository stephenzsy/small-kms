package bootstrap

import (
	"context"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	log "github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
)

type dockerRegistryAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getDockerClient() *dockerclient.Client {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		panic(err)
	}
	return cli
}

func dockerPullImage(ctx context.Context, imageRef string, creds azcore.TokenCredential, tenantID string) error {
	registryLoginUrl, err := acr.ExtractACRLoginServer(imageRef)
	if err != nil {
		return err
	}

	log.Ctx(ctx).Debug().Msgf("Registry login url: %s", registryLoginUrl)
	dcli := getDockerClient()

	authProvider := acr.NewDockerRegistryAuthProvider(registryLoginUrl, creds, tenantID)
	registryAuth, err := authProvider.GetRegistryAuth(ctx)
	if err != nil {
		return err
	}

	out, err := dcli.ImagePull(context.Background(), imageRef, types.ImagePullOptions{
		RegistryAuth: registryAuth,
	})
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	return err
}
