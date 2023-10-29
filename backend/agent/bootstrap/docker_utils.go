package bootstrap

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/google/uuid"
	log "github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/cloudutils"
	"github.com/stephenzsy/small-kms/backend/internal/tokenutils/acr"
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
	registryLoginUrl, err := cloudutils.ExtractACRLoginServer(imageRef)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Registry login url: %s", registryLoginUrl)

	registryEndpoint := "https://" + registryLoginUrl

	acrAuthCli := acr.NewAuthenticationClient(registryEndpoint, creds, &acr.AuthenticationClientOptions{
		TenantID: tenantID,
	})
	token, err := acrAuthCli.ExchagneAADTokenForACRRefreshToken(ctx, registryLoginUrl)
	if err != nil {
		return fmt.Errorf("failed to exchange token: %w", err)
	}

	dcli := getDockerClient()
	dra := dockerRegistryAuth{
		Username: uuid.Nil.String(),
		Password: *token.RefreshToken,
	}
	dockerRegistryAuthJson, err := json.Marshal(dra)
	if err != nil {
		return err
	}

	out, err := dcli.ImagePull(context.Background(), imageRef, types.ImagePullOptions{
		RegistryAuth: base64.RawURLEncoding.EncodeToString(dockerRegistryAuthJson),
	})
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	return err
}
