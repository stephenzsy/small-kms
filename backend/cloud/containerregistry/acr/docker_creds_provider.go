package acr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type acrDockerRegistryAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`

	expiresAt  jwt.NumericDate
	authString string
}

type DockerRegistryAuthProvider struct {
	loginServer         string
	currentRegistryAuth *acrDockerRegistryAuth
	acrAuthCli          *authenticationClient
	jwtParser           *jwt.Parser
}

func (p *DockerRegistryAuthProvider) GetRegistryAuth(c context.Context) (string, error) {
	if p.currentRegistryAuth == nil || p.currentRegistryAuth.expiresAt.Before(time.Now()) {
		token, err := p.acrAuthCli.ExchagneAADTokenForACRRefreshToken(c, p.loginServer)
		if err != nil {
			return "", err
		}
		claims := jwt.RegisteredClaims{}
		if _, _, err := p.jwtParser.ParseUnverified(*token.RefreshToken, &claims); err != nil {
			return "", err
		}
		p.currentRegistryAuth = &acrDockerRegistryAuth{
			Username:  uuid.UUID{}.String(),
			Password:  *token.RefreshToken,
			expiresAt: *claims.ExpiresAt,
		}
		dockerRegistryAuthJson, err := json.Marshal(&p.currentRegistryAuth)
		if err != nil {
			return "", err
		}
		p.currentRegistryAuth.authString = base64.RawURLEncoding.EncodeToString(dockerRegistryAuthJson)

	}
	return p.currentRegistryAuth.authString, nil
}

func NewDockerRegistryAuthProvider(loginServer string, creds azcore.TokenCredential, tenantID string) *DockerRegistryAuthProvider {
	return &DockerRegistryAuthProvider{
		loginServer: loginServer,
		acrAuthCli: NewAuthenticationClient("https://"+loginServer, creds, &AuthenticationClientOptions{
			TenantID: tenantID,
		}),
		jwtParser: jwt.NewParser(jwt.WithoutClaimsValidation()),
	}
}
