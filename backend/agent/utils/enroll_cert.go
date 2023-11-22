package agentutils

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	agentclient "github.com/stephenzsy/small-kms/backend/agent/client/v2"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/cryptoprovider"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
)

func EnrollCertificate(c context.Context,
	client agentclient.ClientWithResponsesInterface,
	certPolicyID string,
	openFile func(*certmodels.Certificate) (*os.File, error),
	onBehalfOf bool) (*certmodels.Certificate, crypto.PrivateKey, error) {

	bad := func(err error) (*certmodels.Certificate, crypto.PrivateKey, error) {
		return nil, nil, err
	}

	// create keypair
	cryptoStore, err := cryptoprovider.NewCryptoProvider()
	if err != nil {
		return bad(err)
	}
	if cryptoStore == nil {
		return bad(nil)
	}

	privateKey, err := cryptoStore.GenerateRSAKeyPair(2048)
	if err != nil {
		return bad(err)
	}

	publicJwk, err := cloudkey.NewJsonWebKeyFromPublicKey(privateKey.Public())
	if err != nil {
		return bad(err)
	}

	resp, err := client.EnrollCertificateWithResponse(c, models.NamespaceProviderServicePrincipal,
		"me",
		certPolicyID,
		&agentclient.EnrollCertificateParams{
			OnBehalfOfApplication: &onBehalfOf,
		},
		certmodels.EnrollCertificateRequest{
			PublicKey: cloudkey.JsonWebSignatureKey{
				JsonWebKeyBase: *publicJwk,
			},
		})
	if err != nil {
		return bad(err)
	} else if resp.StatusCode() != http.StatusCreated {
		return bad(fmt.Errorf("unexpected status code: %d", resp.StatusCode()))
	}

	pkBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return bad(err)
	}

	certFile, err := openFile(resp.JSON201)
	if err != nil {
		return bad(err)
	}
	defer certFile.Close()

	pem.Encode(certFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkBytes,
	})

	for _, cert := range resp.JSON201.Jwk.CertificateChain {
		pem.Encode(certFile, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		})
	}

	return resp.JSON201, privateKey, err
}
