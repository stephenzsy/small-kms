package admin

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type certificateSigner struct {
	certificateChainPemBlob []byte
	certificate             *x509.Certificate
	privateKey              *keyVaultSigner
	certCUID                kmsdoc.KmsDocID
	certPubKey              crypto.PublicKey
	createAzCertificateResp *azcertificates.CreateCertificateResponse
	rootCAKeyBundle         *azkeys.KeyBundle
	namespaceID             uuid.UUID
}

func (b *certificateSigner) SignatureAlgorithm() x509.SignatureAlgorithm {
	switch b.privateKey.sigAlg {
	case azkeys.SignatureAlgorithmRS256:
		return x509.SHA256WithRSA
	case azkeys.SignatureAlgorithmRS384:
		return x509.SHA384WithRSA
	case azkeys.SignatureAlgorithmRS512:
		return x509.SHA512WithRSA
	case azkeys.SignatureAlgorithmES256:
		return x509.ECDSAWithSHA256
	case azkeys.SignatureAlgorithmES384:
		return x509.ECDSAWithSHA384
	}
	return x509.UnknownSignatureAlgorithm
}

func (s *adminServer) loadCertSigner(ctx context.Context, nsID uuid.UUID,
	tdoc *CertificateTemplateDoc, cert *x509.Certificate) (*certificateSigner, error) {
	signer := certificateSigner{}
	signer.namespaceID = tdoc.IssuerNamespaceID
	if tdoc.NamespaceID == tdoc.IssuerNamespaceID {
		// root ca will create keys in key vault
		keyBundle, err := createAzKey(ctx, s.AzKeysClient(), false, tdoc.KeyProperties, tdoc.KeyStorePath, cert.NotAfter)
		if err != nil {
			return nil, err
		}
		signer.privateKey, err = newKeyVaultSigner(ctx, s.AzKeysClient(), keyBundle.Key,
			tdoc.KeyProperties.Alg.ToAzKeysSignatureAlgorithm())
		if err != nil {
			return nil, err
		}
		signer.certificate = cert
		signer.certPubKey = signer.privateKey.publicKey
		signer.rootCAKeyBundle = &keyBundle

	} else {
		// read certficate doc
		certDoc, err := s.readCertDoc(ctx, tdoc.IssuerNamespaceID, kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForTemplate, tdoc.IssuerTemplateID.GetUUID()))
		if err != nil {
			return nil, err
		}
		// load certificate chain
		signer.certificateChainPemBlob, err = certDoc.fetchCertificatePEMBlob(ctx, s.azBlobContainerClient)
		if err != nil {
			return nil, err
		}
		certBytes, _ := pem.Decode(signer.certificateChainPemBlob)
		signer.certificate, err = x509.ParseCertificate(certBytes.Bytes)
		signer.certCUID = certDoc.GetCUID()
		if err != nil {
			return nil, err
		}

		if tdoc.KeyStorePath != nil && len(*tdoc.KeyStorePath) > 0 {

			// load private key from key store
			if certDoc.KeyInfo.KeyID == nil {
				return nil, fmt.Errorf("issuer certificate %s does not have key ID", signer.certCUID.String())
			}
			signerKID := azkeys.ID(*certDoc.KeyInfo.KeyID)

			keyBundle, err := s.AzKeysClient().GetKey(ctx, signerKID.Name(), signerKID.Version(), nil)
			if err != nil {
				return nil, err
			}
			signer.privateKey, err = newKeyVaultSigner(ctx, s.AzKeysClient(), keyBundle.Key,
				tdoc.KeyProperties.Alg.ToAzKeysSignatureAlgorithm())
			if err != nil {
				return nil, err
			}

			// use create certificate to create managed key in key vault
			azCertResp, err := tdoc.createAzCertificate(ctx, s.AzCertificatesClient(), nsID, cert.Subject.String())
			if err != nil {
				return nil, err
			}
			signer.createAzCertificateResp = &azCertResp
			csr, err := x509.ParseCertificateRequest(azCertResp.CSR)
			if err != nil {
				return nil, err
			}
			signer.certPubKey = csr.PublicKey
		} else {
			// load public key from certificate chain
			return nil, errors.New("not implemented for non key vault distribution")
		}
	}
	return &signer, nil
}

func prepareUnsignedCertificateFromTemplate(
	nsID uuid.UUID, tmplProcessor *certTemplateProcessor) (*x509.Certificate, uuid.UUID, error) {
	// prep certificate
	certID, err := uuid.NewRandom()
	if err != nil {
		return nil, certID, err
	}
	certSerial := big.Int{}
	certSerial.SetBytes(certID[:])
	now := time.Now()
	c := x509.Certificate{
		SerialNumber: &certSerial,
		NotBefore:    now,
	}
	tmplProcessor.processTemplate(&c)

	if err != nil {
		return nil, certID, err
	}
	if isAllowedRootCaNamespace(nsID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLen = 1
		c.MaxPathLenZero = false
		c.BasicConstraintsValid = true
		c.ExtKeyUsage = nil
	} else if isAllowedCaNamespace(nsID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLenZero = true
		c.BasicConstraintsValid = true
		c.ExtKeyUsage = nil
	} else {
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
	}

	return &c, certID, err
}

func certificateSubjectAlternativeNamesToDoc(cert *x509.Certificate) (r *CertificateSubjectAlternativeNames) {
	for _, s := range cert.EmailAddresses {
		if r == nil {
			r = new(CertificateSubjectAlternativeNames)
		}
		r.EmailAddresses = append(r.EmailAddresses, s)
	}
	for _, s := range cert.URIs {

		if r == nil {
			r = new(CertificateSubjectAlternativeNames)
		}
		r.URIs = append(r.URIs, s.String())
	}
	return r
}

// (nsType/nsID) must be verified prior to calling this function
func (s *adminServer) createCertificateFromTemplate(ctx context.Context, nsID uuid.UUID,
	t *certTemplateProcessor) (*CertDoc, []byte, error) {
	c, certID, err := prepareUnsignedCertificateFromTemplate(nsID, t)
	if err != nil {
		return nil, nil, err
	}
	return s.createCertificateFromTemplateWithCert(ctx, nsID, t.tmplDoc, c, certID)
}

func (s *adminServer) createCertificateFromTemplateWithCert(ctx context.Context, nsID uuid.UUID,
	t *CertificateTemplateDoc, c *x509.Certificate, certID uuid.UUID) (*CertDoc, []byte, error) {

	// prep signer
	signer, err := s.loadCertSigner(ctx, nsID, t, c)
	if err != nil {
		return nil, nil, err
	}
	c.SignatureAlgorithm = signer.SignatureAlgorithm()
	log.Info().Msgf("signer %s", signer.certCUID.String())
	// Sign cert
	certSigned, err := x509.CreateCertificate(nil, c, signer.certificate, signer.certPubKey, signer.privateKey)
	if err != nil {
		return nil, nil, err
	}
	certParsed, err := x509.ParseCertificate(certSigned)
	if err != nil {
		return nil, nil, err
	}
	log.Info().Msgf("Certificate signed and validated, prepare to upload: %s", certID)

	// encode to pem
	pemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certSigned,
	}
	bb := bytes.Buffer{}
	err = pem.Encode(&bb, &pemBlock)
	if err != nil {
		return nil, nil, err
	}
	if len(signer.certificateChainPemBlob) > 0 {
		// attach chain
		_, err = bb.Write(signer.certificateChainPemBlob)
		if err != nil {
			return nil, nil, err
		}
	}

	// finalize certificate on keyvault if there is an outstanding request
	var mergeCertificateResponse *azcertificates.MergeCertificateResponse
	if signer.createAzCertificateResp != nil {
		mergeRequestCerts := make([][]byte, 1, 3)
		mergeRequestCerts[0] = certSigned
		pemBlob := signer.certificateChainPemBlob
		for block, rest := pem.Decode(pemBlob); block != nil; block, rest = pem.Decode(rest) {
			mergeRequestCerts = append(mergeRequestCerts, block.Bytes)
		}

		azcMcResp, err := s.AzCertificatesClient().MergeCertificate(ctx, signer.createAzCertificateResp.ID.Name(), azcertificates.MergeCertificateParameters{
			X509Certificates: mergeRequestCerts,
		}, nil)
		if err != nil {
			return nil, nil, err
		}
		mergeCertificateResponse = &azcMcResp
	}

	blobName := fmt.Sprintf("%s/%s.pem", nsID, certID)
	// prepare certDocument
	certDoc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, certID),
			NamespaceID: nsID,
		},
		TemplateID:              t.ID,
		Subject:                 certParsed.Subject.String(),
		SubjectBase:             t.Subject.String(),
		NotBefore:               c.NotBefore,
		NotAfter:                c.NotAfter,
		SubjectAlternativeNames: certificateSubjectAlternativeNamesToDoc(certParsed),
		CertStorePath:           blobName, // certificate storage path in blob storage
		IssuerNamespaceID:       signer.namespaceID,
		IssuerCertificateID:     signer.certCUID,
		CommonName:              certParsed.Subject.CommonName,
	}
	certDoc.KeyInfo.populateBriefFromCertificate(certParsed)
	certDoc.FingerprintSHA1Hex = base64UrlToHexStr(*certDoc.KeyInfo.CertificateThumbprint)
	// populate x5u
	if mergeCertificateResponse != nil {
		certDoc.KeyInfo.KeyID = (*string)(mergeCertificateResponse.KID)
		certDoc.KeyInfo.CertificateURL = (*string)(mergeCertificateResponse.ID)
	} else if signer.rootCAKeyBundle != nil {
		certDoc.KeyInfo.KeyID = (*string)(signer.rootCAKeyBundle.Key.KID)
	}

	// upload to blob
	certChainPemBlob := bb.Bytes()
	blobUrl, err := certDoc.storeCertificatePEMBlob(ctx, s.azBlobContainerClient, certChainPemBlob)
	if err != nil {
		return nil, nil, err
	}
	if certDoc.KeyInfo.CertificateURL == nil {
		certDoc.KeyInfo.CertificateURL = blobUrl
	}

	// upload to blob
	return &certDoc, certChainPemBlob, nil
}
