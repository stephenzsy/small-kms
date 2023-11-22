package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"

	agentclient "github.com/stephenzsy/small-kms/backend/agent/client/v2"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
)

type ServicePrincipalBootstraper struct {
}

func NewServicePrincipalBootstraper() *ServicePrincipalBootstraper {
	return &ServicePrincipalBootstraper{}
}

func (*ServicePrincipalBootstraper) Bootstrap(c context.Context, certPolicyID string, certPath string, tokenCacheFile string) error {
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
	var ok bool
	envSvc := common.NewEnvService()
	if baseUrl, ok = envSvc.Require(agentcommon.EnvKeyAPIBaseURL, common.IdentityEnvVarPrefixApp); !ok {
		return envSvc.ErrMissing(agentcommon.EnvKeyAPIBaseURL)
	} else if apiAuthScope, ok = envSvc.Require(agentcommon.EnvKeyAPIAuthScope, common.IdentityEnvVarPrefixApp); !ok {
		return envSvc.ErrMissing(agentcommon.EnvKeyAPIAuthScope)
	}

	appTokenCache := newAppTokenCache(tokenCacheFile)
	pubClient, authResult, err := getAppWithSharedTokenCache(c, appTokenCache, true, false)
	if err != nil {
		return err
	}

	// nbf := jwt.NewNumericDate(time.Now())

	// t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
	// 	Audience:  jwt.ClaimStrings{"00000003-0000-0000-c000-000000000000"},
	// 	NotBefore: nbf,
	// 	ExpiresAt: jwt.NewNumericDate(nbf.Time.Add(10 * time.Minute)),
	// 	Issuer:    string(namespaceIdentifier),
	// })
	// signedToken, err := t.SignedString(privateKey)
	// if err != nil {
	// 	return err
	// }

	client, err := agentclient.NewClientWithResponses(baseUrl,
		agentclient.WithRequestEditorFn(common.ToSilenTokenRequestEditorFn(pubClient, apiAuthScope, authResult.Account)))
	if err != nil {
		return err
	}

	_, _, err = agentutils.EnrollCertificate(c, client, certPolicyID,
		func(_ *certmodels.Certificate) (*os.File, error) {
			return os.OpenFile(certPath, os.O_CREATE|os.O_WRONLY, 0400)
		}, true)
	return err
}
