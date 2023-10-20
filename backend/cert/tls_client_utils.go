package cert

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func CertDocKeyPair(c context.Context, certDoc *CertDoc) (tls.Certificate, error) {
	fail := func(err error) (tls.Certificate, error) { return tls.Certificate{}, err }

	certPEMBlock, err := certDoc.FetchCertificatePEMBlob(c)
	if err != nil {
		return fail(err)
	}

	var cert tls.Certificate
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}

	if len(cert.Certificate) == 0 {
		return fail(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate"))
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fail(err)
	}

	cert.PrivateKey = newKeyVaultSignerWithExistingPublicKey(c,
		utils.ToPtr(azkeys.ID(certDoc.CertSpec.KID)),
		x509Cert.PublicKey,
		certDoc.CertSpec.Alg.ToAzKeysSignatureAlgorithm())

	return cert, nil
}
