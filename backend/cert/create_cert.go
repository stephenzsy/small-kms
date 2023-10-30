package cert

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type existingPublicKeyProvider struct {
	crypto.PublicKey
}

// Cleanup implements kv.AzCertCSRProvider.
func (*existingPublicKeyProvider) Cleanup(context.Context) {
	// do nothing
}

// CollectCerts implements kv.AzCertCSRProvider.
func (*existingPublicKeyProvider) CollectCerts(context.Context, [][]byte) (*azcertificates.MergeCertificateResponse, error) {
	return nil, nil
}

// GetCSRPublicKey implements kv.AzCertCSRProvider.
func (p *existingPublicKeyProvider) GetCSRPublicKey(context.Context) (crypto.PublicKey, error) {
	return p.PublicKey, nil
}

var _ kv.AzCertCSRProvider = (*existingPublicKeyProvider)(nil)

func signCertificate(
	c context.Context,
	template, parent *x509.Certificate,
	csrProvider kv.AzCertCSRProvider, signer kv.AzCertSigner,
	signerLocator base.DocFullIdentifier, signingCertChain [][]byte) (*CertDocSigningPatch, error) {
	err := signer.Load(c)
	if err != nil {
		return nil, err
	}
	defer csrProvider.Cleanup(c)

	csrPubKey, err := csrProvider.GetCSRPublicKey(c)
	if err != nil {
		return nil, err
	}
	signedCert, err := x509.CreateCertificate(nil, template, parent, csrPubKey, signer)
	if err != nil {
		return nil, err
	}
	collectCert := make([][]byte, 0, len(signingCertChain)+1)
	collectCert = append(collectCert, signedCert)
	collectCert = append(collectCert, signingCertChain...)
	mergeResp, err := csrProvider.CollectCerts(c, collectCert)
	if err != nil {
		return nil, err
	}

	patch := new(CertDocSigningPatch)
	// patch.KeySpec.PopulatePublicKey(csrPubKey)
	patch.KeySpec.CertificateChain = utils.MapSlice(collectCert, func(certBytes []byte) base.Base64RawURLEncodedBytes {
		return base.Base64RawURLEncodedBytes(certBytes)
	})
	if mergeResp != nil {
		patch.KeySpec.KeyID = utils.ToPtr(string(*mergeResp.KID))
		patch.KeyVaultStore = CertDocKeyVaultStore{
			Name: mergeResp.ID.Name(),
			ID:   string(*mergeResp.ID),
			SID:  string(*mergeResp.SID),
		}
	} else {
		privateLabel := "_private"
		patch.KeySpec.KeyID = &privateLabel
	}
	certSHA1 := sha1.Sum(signedCert)
	patch.KeySpec.X5t = base.Base64RawURLEncodedBytes(certSHA1[:])
	certSHA256 := sha256.Sum256(signedCert)
	patch.KeySpec.X5tS256 = base.Base64RawURLEncodedBytes(certSHA256[:])

	patch.Issuer = signerLocator
	return patch, nil
}

func createCertFromPolicy(c context.Context, policyRID Identifier, publicKey crypto.PublicKey) (*CertDoc, error) {
	policyDoc, err := ReadCertPolicyDoc(c, policyRID)
	if err != nil {
		return nil, err
	}

	doc := new(CertDoc)
	nsCtx := ns.GetNSContext(c)
	doc.Init(nsCtx.Kind(), nsCtx.Identifier(), policyDoc)
	docService := base.GetAzCosmosCRUDService(c)

	switch nsCtx.Kind() {
	case base.NamespaceKindRootCA:
		c = ctx.Elevate(c)
		signer := kv.NewAzCertSelfSigner(doc.getCSRProviderParams(), doc.getSigningParams())
		cert := doc.getX509CertTemplate()
		cert.SignatureAlgorithm = doc.getX509SignatureAlgorithm()

		err = docService.Create(c, doc, nil)
		if err != nil {
			return nil, err
		}

		patch, err := signCertificate(c, cert, cert, signer, signer, doc.GetStorageFullIdentifier(), nil)
		if err != nil {
			return nil, err
		}
		err = doc.applyPatch(c, docService, patch)
		if err != nil {
			return nil, err
		}
	case base.NamespaceKindIntermediateCA,
		base.NamespaceKindServicePrincipal:
		c = ctx.Elevate(c)
		// load certDoc of signer
		issuerNamespace := policyDoc.IssuerNamespace
		if issuerNamespace.Kind() == nsCtx.Kind() && issuerNamespace.Identifier() == nsCtx.Identifier() {
			return nil, fmt.Errorf("%w: this operation does not support creating self-signed certificate", base.ErrResponseStatusBadRequest)
		}

		signerDoc, err := getNamespaceIssuerCert(c, issuerNamespace)
		if err != nil {
			return nil, err
		}
		signerCert, signerChain, err := signerDoc.getIssuedX509Certificate()
		if err != nil {
			return nil, err
		}
		signer := kv.NewAzCertSigner(signerDoc.getSigningParams(), signerCert.PublicKey)
		var csrProvider kv.AzCertCSRProvider
		if publicKey != nil {
			csrProvider = &existingPublicKeyProvider{publicKey}
		} else {
			csrProvider = kv.NewAzCSRProvider(doc.getCSRProviderParams())
		}
		cert := doc.getX509CertTemplate()
		cert.SignatureAlgorithm = signerDoc.getX509SignatureAlgorithm()

		err = docService.Create(c, doc, nil)
		if err != nil {
			return nil, err
		}

		patch, err := signCertificate(c, cert, signerCert, csrProvider, signer, doc.GetStorageFullIdentifier(), signerChain)
		if err != nil {
			return nil, err
		}
		err = doc.applyPatch(c, docService, patch)
		if err != nil {
			return nil, err
		}

	}

	return doc, err
}
