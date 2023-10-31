package bootstrap

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
	"github.com/stephenzsy/small-kms/backend/key"
)

type ServicePrincipalBootstraper struct {
}

func NewServicePrincipalBootstraper() *ServicePrincipalBootstraper {
	return &ServicePrincipalBootstraper{}
}

func (*ServicePrincipalBootstraper) Bootstrap(c context.Context, namespaceIdentifier base.Identifier, certPolicyIdentifer base.Identifier, certPath string, tokenCacheFile string) error {
	if certPath == "" {
		return errors.New("missing client cert path")
	}
	if _, err := os.Stat(certPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		fmt.Println("client cert already exists, skipping")
		return nil
	}

	// create keypair
	cryptoStore, err := cryptoprovider.NewCryptoProvider()
	if err != nil {
		return err
	}
	if cryptoStore == nil {
		return nil
	}

	var baseUrl, apiAuthScope string
	if baseUrl = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "API_BASE_URL", ""); baseUrl == "" {
		return errors.New("missing API_URL_BASE")
	} else if apiAuthScope = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "API_AUTH_SCOPE", ""); apiAuthScope == "" {
		return errors.New("missing APP_API_AUTH_SCOPE")
	}

	appTokenCache := newAppTokenCache(tokenCacheFile)
	pubClient, authResult, err := getAppWithSharedTokenCache(c, appTokenCache, true, false)
	if err != nil {
		return err
	}

	privateKey, err := cryptoStore.GenerateRSAKeyPair(2048)
	if err != nil {
		return err
	}

	nbf := jwt.NewNumericDate(time.Now())

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"00000003-0000-0000-c000-000000000000"},
		NotBefore: nbf,
		ExpiresAt: jwt.NewNumericDate(nbf.Time.Add(10 * time.Minute)),
		Issuer:    namespaceIdentifier.String(),
	})
	signedToken, err := t.SignedString(privateKey)
	if err != nil {
		return err
	}

	client, err := agentclient.NewClientWithResponses(baseUrl,
		agentclient.WithRequestEditorFn(common.ToSilenTokenRequestEditorFn(pubClient, apiAuthScope, authResult.Account)))
	if err != nil {
		return err
	}

	resp, err := client.EnrollCertificateWithResponse(c, base.NamespaceKindServicePrincipal,
		namespaceIdentifier,
		certPolicyIdentifer,
		agentclient.EnrollCertificateRequest{
			PublicKey: toJwk(privateKey.Public()),
			Proof:     signedToken,
		})
	if err != nil {
		return err
	}

	pkBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	certFile, err := os.OpenFile(certPath, os.O_CREATE|os.O_WRONLY, 0400)
	if err != nil {
		return err
	}
	defer certFile.Close()

	pem.Encode(certFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkBytes,
	})

	for _, cert := range resp.JSON200.Jwk.CertificateChain {
		pem.Encode(certFile, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		})
	}

	return nil
}

func toJwk(k crypto.PublicKey) (jwk key.JsonWebKey) {
	if rsaPubKey, ok := k.(*rsa.PublicKey); ok {
		jwk.Kty = key.JsonWebKeyTypeRSA
		jwk.N = base.Base64RawURLEncodedBytes(rsaPubKey.N.Bytes())
		jwk.E = base.Base64RawURLEncodedBytes(big.NewInt(int64(rsaPubKey.E)).Bytes())
	}
	return
}
