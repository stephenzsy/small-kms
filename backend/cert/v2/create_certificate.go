package cert

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func signCertificate(
	c context.Context,
	template, parent *x509.Certificate,
	csrProvider kv.AzCertCSRProvider, signer kv.AzCertSigner,
	signerLocator base.SLocator, signingCertChain [][]byte) (*CertDocSigningPatch, error) {
	err := signer.Load(c)
	if err != nil {
		return nil, err
	}
	createCertResp, err := csrProvider.GetCSR(c)
	if err != nil {
		return nil, err
	}
	defer csrProvider.Cleanup(c, createCertResp)
	csrParsed, err := x509.ParseCertificateRequest(createCertResp.CSR)
	if err != nil {
		return nil, err
	}
	signedCert, err := x509.CreateCertificate(nil, template, parent, csrParsed.PublicKey, signer)
	if err != nil {
		return nil, err
	}
	collectCert := make([][]byte, 0, len(signingCertChain)+1)
	collectCert = append(collectCert, signedCert)
	collectCert = append(collectCert, signingCertChain...)
	mergeResp, err := csrProvider.CollectCerts(c, createCertResp, collectCert)
	if err != nil {
		return nil, err
	}

	patch := new(CertDocSigningPatch)
	patch.KeySpec.CertificateChain = utils.MapSlice(collectCert, func(certBytes []byte) base.Base64RawURLEncodedBytes {
		return base.Base64RawURLEncodedBytes(certBytes)
	})
	patch.KeySpec.KeyID = utils.ToPtr(string(*mergeResp.KID))
	certSHA1 := sha1.Sum(signedCert)
	patch.KeySpec.X5t = utils.ToPtr(base.Base64RawURLEncodedBytes(certSHA1[:]))
	certSHA256 := sha256.Sum256(signedCert)
	patch.KeySpec.X5tS256 = utils.ToPtr(base.Base64RawURLEncodedBytes(certSHA256[:]))

	patch.KeyVaultStore = CertDocKeyVaultStore{
		Name: mergeResp.ID.Name(),
		ID:   string(*mergeResp.ID),
		SID:  string(*mergeResp.SID),
	}

	patch.Issuer = signerLocator
	return patch, nil
}

func createCertFromPolicy(c context.Context, policyRID Identifier) (*CertDoc, error) {

	policyDoc, err := getCertPolicy(c, policyRID)
	if err != nil {
		return nil, err
	}

	doc := new(CertDoc)
	nsCtx := ns.GetNSContext(c)
	doc.Init(nsCtx.Kind(), nsCtx.Identifier(), policyDoc)
	docService := base.GetAzCosmosCRUDService(c)
	err = docService.Create(c, doc, nil)
	if err != nil {
		return nil, err
	}

	if nsCtx.Kind() == base.NamespaceKindRootCA {
		c = ctx.Elevate(c)
		sp := doc.getSigningParams()
		cert := doc.getX509CertTemplate()
		signer, err := kv.NewAzCertSelfSigner(sp)
		if err != nil {
			return nil, err
		}
		patch, err := signCertificate(c, cert, cert, signer, signer, doc.GetSLocator(), nil)
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
