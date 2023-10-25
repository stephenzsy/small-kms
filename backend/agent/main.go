package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/agent/agentserver"
	cm "github.com/stephenzsy/small-kms/backend/agent/configmanager"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/tokenutils/acr"
	"github.com/stephenzsy/small-kms/backend/shared"
)

const DefaultEnvVarTenantID = "AZURE_TENANT_ID"
const DefaultEnvVarClientID = "AZURE_CLIENT_ID"
const DefaultEnvVarCertBundlePath = "AZURE_CERT_BUNDLE"
const DefaultEnvVarApiBaseUrl = "SMALLKMS_API_BASE_URL"
const DefaultEnvVarApiScope = "SMALLKMS_API_SCOPE"

var BuildID = "dev"

func getDockerClient() *dockerclient.Client {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion("1.43"))
	if err != nil {
		panic(err)
	}
	return cli
}

type dockerRegistryAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func dockerPull(ctx context.Context, creds azcore.TokenCredential, tenantID, imageRef string) {
	parsedUrl, err := url.Parse("https://" + imageRef)
	if err != nil {
		log.Panicf("Failed to parse image ref: %v\n", err)
	}
	registryLoginUrl := parsedUrl.Host
	log.Printf("Registry login url: %s\n", registryLoginUrl)

	registryEndpoint := "https://" + registryLoginUrl

	acrAuthCli := acr.NewAuthenticationClient(registryEndpoint, creds, &acr.AuthenticationClientOptions{
		TenantID: tenantID,
	})
	token, err := acrAuthCli.ExchagneAADTokenForACRRefreshToken(ctx, registryLoginUrl)
	if err != nil {
		log.Panicf("Failed to exchange token: %v\n", err)
	}

	dcli := getDockerClient()
	dra := dockerRegistryAuth{
		Username: uuid.Nil.String(),
		Password: *token.RefreshToken,
	}
	dockerRegistryAuthJson, err := json.Marshal(dra)
	if err != nil {
		panic(err)
	}

	out, err := dcli.ImagePull(context.Background(), imageRef, types.ImagePullOptions{
		RegistryAuth: base64.RawURLEncoding.EncodeToString(dockerRegistryAuthJson),
	})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

var myNamespacIdentifier = shared.StringIdentifier("me")

func bootstrapActiveHost(skipDockerPullPtr bool) {

	var clientID, tenantID, bundlePath, apiBaseUrl, apiScope string
	var ok bool
	if tenantID, ok = os.LookupEnv(DefaultEnvVarTenantID); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarTenantID)
	}
	if clientID, ok = os.LookupEnv(DefaultEnvVarClientID); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarClientID)
	}
	if bundlePath, ok = os.LookupEnv(DefaultEnvVarCertBundlePath); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarCertBundlePath)
	}
	if apiBaseUrl, ok = os.LookupEnv(DefaultEnvVarApiBaseUrl); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarApiBaseUrl)
	}
	if apiScope, ok = os.LookupEnv(DefaultEnvVarApiScope); !ok {
		log.Panicf("Environment variable %s is not set\n", DefaultEnvVarApiScope)
	}
	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		log.Panicf("Failed to read certificate bundle: %v\n", err)
	}

	var privateKey crypto.PrivateKey
	var x509Certs []*x509.Certificate
	for block, rest := pem.Decode(bundleBytes); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "PRIVATE KEY" {
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				{
					log.Panicf("Failed to parse private key: %v\n", err)
				}
			}
		} else if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Panicf("Failed to parse certificate: %v\n", err)
			}
			x509Certs = append(x509Certs, cert)
			break
		}

	}

	ctx := context.Background()

	creds, err := azidentity.NewClientCertificateCredential(tenantID, clientID, x509Certs, privateKey, nil)
	if err != nil {
		log.Panicf("Failed to create credential: %v\n", err)
	}
	client, err := agentclient.NewClientWithCreds(apiBaseUrl, creds, []string{apiScope}, tenantID)
	if err != nil {
		log.Panicf("Failed to create client: %v\n", err)
	}
	resp, err := client.GetAgentConfigurationWithResponse(ctx, shared.NamespaceKindServicePrincipal, myNamespacIdentifier, shared.AgentConfigNameActiveHostBootstrap, nil)
	if err != nil {
		log.Panicf("Failed to check in: %v\n", err)
	}

	config, err := resp.JSON200.Config.AsAgentConfigurationAgentActiveHostBootstrap()
	if err != nil {
		log.Panicf("Failed to parse config: %v\n", err)
	}

	if !skipDockerPullPtr {
		dockerPull(ctx, creds, tenantID, config.ControllerContainer.ImageRefStr)
	}
}

func _oldmain() {
	// Find .env file
	skipDockerPullPtr := flag.Bool("skip-docker-pull", false, "skip docker pull")
	envFilePathPtr := flag.String("env", "", "path to .env file")
	slotPtr := flag.Uint("slot", 0, "slot")
	//	skipTlsPtr := flag.Bool("skip-tls", false, "skip tls")
	flag.Parse()

	if *envFilePathPtr != "" {
		err := godotenv.Load(*envFilePathPtr)
		if err != nil {
			log.Printf("Error loading environment file: %s: %s\n", *envFilePathPtr, err.Error())
		}
	}

	configDir := common.MustGetenv("AGENT_CONFIG_DIR")

	args := flag.Args()
	if len(args) >= 2 {
		switch args[0] {
		case "bootstrap":
			switch args[1] {
			case "active-host":
				bootstrapActiveHost(*skipDockerPullPtr)
				return
			}
		case "server":
			e := echo.New()
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			e.TLSServer.Addr = args[1]
			configManager, err := cm.NewConfigManager(BuildID, configDir, args[1], uint32(*slotPtr))
			if err != nil {
				log.Panicf("Failed to create config manager: %v\n", err)
			}
			s, err := agentserver.NewServer()
			if err != nil {
				log.Panicf("Failed to create server: %v\n", err)
			}
			agentserver.RegisterHandlers(e, s)

			configManager.Manage(e)

			cm.StartConfigManagerWithGracefulShutdown(context.Background(), configManager)
			//bootstrapServer(args[1], *skipTlsPtr)
			return
		}
	}

	fmt.Printf("Version: %s\n", BuildID)
	fmt.Printf("Usage: %s bootstrap active-host\n", os.Args[0])
	fmt.Printf("Usage: %s server :8443\n", os.Args[0])

}
