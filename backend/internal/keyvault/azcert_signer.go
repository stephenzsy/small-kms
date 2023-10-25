package kv

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/rs/zerolog/log"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type AzCertSigner interface {
	crypto.Signer
	Load(context.Context) error
}

type AzCertCSRProvider interface {
	GetCSRPublicKey(context.Context) (crypto.PublicKey, error)
	CollectCerts(context.Context, [][]byte) (*azcertificates.MergeCertificateResponse, error)
	Cleanup(context.Context)
}

type CSRProviderParams struct {
	CertName      string
	KeyProperties azcertificates.KeyProperties
}

type SigningParams struct {
	CertID azcertificates.ID
	SigAlg azkeys.SignatureAlgorithm
}

type azcertKeyPair struct {
	signingCtx    *context.Context
	csrParams     *CSRProviderParams
	signingParams *SigningParams
	isSelfSigning bool

	temporalCertID *azcertificates.ID
	certPublicKey  crypto.PublicKey
	skipCleanup    bool
}

// Cleanup implements AzCertCSRProvider.
func (kp *azcertKeyPair) Cleanup(c context.Context) {
	if !kp.skipCleanup {
		client := getAzKeyVaultService(c).AzCertificatesClient()

		_, err := client.DeleteCertificateOperation(c, kp.csrParams.CertName, nil)
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("failed to delete certificate operation")
		}
	}
}

// CollectCerts implements AzCertCSRProvider.
func (kp *azcertKeyPair) CollectCerts(c context.Context, certs [][]byte) (*azcertificates.MergeCertificateResponse, error) {
	client := getAzKeyVaultService(c).AzCertificatesClient()
	resp, err := client.MergeCertificate(c, kp.csrParams.CertName, azcertificates.MergeCertificateParameters{
		X509Certificates: certs,
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Enabled: to.Ptr(true),
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	kp.skipCleanup = true
	if kp.temporalCertID != nil {

		// disable temporal cert
		_, err = client.UpdateCertificate(c, kp.temporalCertID.Name(), kp.temporalCertID.Version(), azcertificates.UpdateCertificateParameters{
			CertificateAttributes: &azcertificates.CertificateAttributes{
				Enabled: to.Ptr(false),
			},
		}, nil)
		if err != nil {
			return nil, err
		}
	}
	return &resp, nil
}

// GetCSR implements AzCertCSRProvider.
func (kp *azcertKeyPair) GetCSRPublicKey(c context.Context) (crypto.PublicKey, error) {
	if kp.signingCtx == nil {
		if err := kp.Load(c); err != nil {
			return nil, err
		}
	}
	client := getAzKeyVaultService(c).AzCertificatesClient()
	subject := pkix.Name{CommonName: "dummy cert"}.String()
	params := azcertificates.CreateCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Enabled: to.Ptr(true),
		},
		CertificatePolicy: &azcertificates.CertificatePolicy{
			X509CertificateProperties: &azcertificates.X509CertificateProperties{
				Subject: &subject,
			},
			KeyProperties: &kp.csrParams.KeyProperties,
			SecretProperties: &azcertificates.SecretProperties{
				ContentType: to.Ptr("application/x-pem-file"),
			},
		},
	}
	if kp.isSelfSigning {
		params.CertificatePolicy.KeyProperties.ReuseKey = to.Ptr(true)
	}
	if resp, err := client.CreateCertificate(c, kp.csrParams.CertName, params, nil); err != nil {
		return nil, err
	} else if csrParsed, err := x509.ParseCertificateRequest(resp.CSR); err != nil {
		return nil, err
	} else {
		return csrParsed.PublicKey, nil
	}
}

// Load implements AzCertSigner.
func (kp *azcertKeyPair) Load(c context.Context) error {
	c = ctx.Elevate(c)

	if kp.isSelfSigning {
		// elevate context to ignore cancellation

		client := getAzKeyVaultService(c).AzCertificatesClient()
		subject := pkix.Name{CommonName: "dummy cert for key"}.String()
		params := azcertificates.CreateCertificateParameters{
			CertificateAttributes: &azcertificates.CertificateAttributes{
				Enabled: to.Ptr(true), // we want to use the key, so must be enabled
			},
			CertificatePolicy: &azcertificates.CertificatePolicy{
				X509CertificateProperties: &azcertificates.X509CertificateProperties{
					Subject: &subject,
				},
				KeyProperties: &kp.csrParams.KeyProperties,
				IssuerParameters: &azcertificates.IssuerParameters{
					Name: to.Ptr("Self"),
				},
				SecretProperties: &azcertificates.SecretProperties{
					ContentType: to.Ptr("application/x-pem-file"),
				},
			},
		}
		resp, err := client.CreateCertificate(c, kp.csrParams.CertName, params, nil)
		if err != nil {
			return err
		}
		certID := resp.ID
		status := resp.Status
		for status != nil && *status == "inProgress" {
			// wait for 5 seconds
			time.Sleep(5 * time.Second)
			resp, err := client.GetCertificateOperation(c, certID.Name(), nil)
			if err != nil {
				return err
			}
			status = resp.Status
		}
		getCertResp, err := client.GetCertificate(c, certID.Name(), "", nil)
		if err != nil {
			return err
		}
		kp.signingParams.CertID = *getCertResp.ID
		kp.temporalCertID = getCertResp.ID
		parsed, err := x509.ParseCertificate(getCertResp.CER)
		if err != nil {
			return err
		}
		kp.certPublicKey = parsed.PublicKey
	}

	kp.signingCtx = &c
	return nil
}

// Public implements AzCertSigner.
func (kp *azcertKeyPair) Public() crypto.PublicKey {
	return kp.certPublicKey
}

// Sign implements AzCertSigner.
func (p *azcertKeyPair) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if p.signingCtx == nil {
		return nil, errors.New("signer not loaded")
	}
	resp, err := getAzKeyVaultService(*p.signingCtx).AzKeysClient().
		Sign(*p.signingCtx, p.signingParams.CertID.Name(), p.signingParams.CertID.Version(), azkeys.SignParameters{
			Value:     digest,
			Algorithm: &p.signingParams.SigAlg,
		}, nil)
	if err != nil {
		return nil, err
	}
	return toX509Signature(resp.Result, p.signingParams.SigAlg)
}

var _ AzCertSigner = (*azcertKeyPair)(nil)
var _ AzCertCSRProvider = (*azcertKeyPair)(nil)

func NewAzCertSelfSigner(pCsr CSRProviderParams, pSigning SigningParams) *azcertKeyPair {
	return &azcertKeyPair{
		signingParams: &pSigning,
		csrParams:     &pCsr,
		isSelfSigning: true,
	}
}

func NewAzCertSigner(pSigning SigningParams, publicKey crypto.PublicKey) AzCertSigner {
	return &azcertKeyPair{
		signingParams: &pSigning,
		isSelfSigning: false,
		certPublicKey: publicKey,
	}
}

func NewAzCSRProvider(pCsr CSRProviderParams) AzCertCSRProvider {
	return &azcertKeyPair{
		csrParams:     &pCsr,
		isSelfSigning: false,
	}
}
