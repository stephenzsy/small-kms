package serviceprincipal

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/internal/certstore"
	wincryptostore "github.com/stephenzsy/small-kms/backend/internal/certstore/windows"
)

type ServicePrincipalBootstraper struct {
}

func NewServicePrincipalBootstraper() *ServicePrincipalBootstraper {
	return &ServicePrincipalBootstraper{}
}

func (*ServicePrincipalBootstraper) Bootstrap(c context.Context, certPath, tokenCacheFile string) error {
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
	cryptoStore, err := getCryptoStoreProvider()
	if err != nil {
		return err
	}
	if cryptoStore == nil {
		return nil
	}

	appTokenCache := newAppTokenCache(tokenCacheFile)
	defer appTokenCache.Close()
	var issuerId string
	_, authResult, err := getAppWithSharedTokenCache(c, appTokenCache, true, false)
	if err != nil {
		issuerId = authResult.Account.LocalAccountID
	}

	ksession, err := cryptoStore.CreateRSAKeySession("smallkms", 2048, false)
	if err != nil {
		return err
	}
	defer ksession.Close()

	nbf := jwt.NewNumericDate(time.Now())

	t := jwt.NewWithClaims(&ksessionSigningMethod{ksession}, jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"00000003-0000-0000-c000-000000000000"},
		NotBefore: nbf,
		ExpiresAt: jwt.NewNumericDate(nbf.Time.Add(10 * time.Minute)),
		Issuer:    issuerId,
	})
	signedToken, err := t.SignedString(ksession)
	if err != nil {
		return err
	}
	log.Info().Msg(signedToken)

	return nil
}

func getCryptoStoreProvider() (certstore.CryptoStoreProvider, error) {
	if runtime.GOOS == "windows" {
		return wincryptostore.NewWindowsNCryptCryptoStoreProvider(), nil
	}
	return nil, nil
}

type ksessionSigningMethod struct {
	ksession certstore.KeySession
}

// Alg implements jwt.SigningMethod.
func (*ksessionSigningMethod) Alg() string {
	return "RS256"
}

// Sign implements jwt.SigningMethod.
func (ksm *ksessionSigningMethod) Sign(signingString string, key interface{}) ([]byte, error) {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(signingString))
	return ksm.ksession.Sign(nil, hasher.Sum(nil), crypto.SHA256.HashFunc())
}

// Verify implements jwt.SigningMethod.
func (*ksessionSigningMethod) Verify(signingString string, sig []byte, key interface{}) error {
	panic("unimplemented")
}

var _ jwt.SigningMethod = (*ksessionSigningMethod)(nil)
