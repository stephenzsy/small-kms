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

func (alg JwkAlg) toAzKeysSignatureAlgorithm() azkeys.SignatureAlgorithm {
	switch alg {
	case AlgRS256:
		return azkeys.SignatureAlgorithmRS256
	case AlgRS512:
		return azkeys.SignatureAlgorithmRS512
	}
	return azkeys.SignatureAlgorithmRS384
}

func (s *adminServer) loadCertSigner(ctx context.Context, nsType NamespaceTypeShortName, nsID uuid.UUID,
	tdoc *CertificateTemplateDoc, cert *x509.Certificate, tmplData *TemplateVarData) (*certificateSigner, error) {
	signer := certificateSigner{}
	if tdoc.NamespaceID == tdoc.IssuerNamespaceID {
		// root ca will create keys in key vault
		keyBundle, err := tdoc.createAzKey(ctx, s.AzKeysClient(), nsType, cert)
		if err != nil {
			return nil, err
		}
		signer.privateKey, err = newKeyVaultSigner(ctx, s.AzKeysClient(), keyBundle.Key,
			tdoc.KeyProperties.Alg.toAzKeysSignatureAlgorithm())
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
				tdoc.KeyProperties.Alg.toAzKeysSignatureAlgorithm())
			if err != nil {
				return nil, err
			}

			// use create certificate to create managed key in key vault
			azCertResp, err := tdoc.createAzCertificate(ctx, s.AzCertificatesClient(), nsType, tmplData)
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

func prepareUnsignedCertificateFromTemplate(nsType NamespaceTypeShortName,
	nsID uuid.UUID, t *CertificateTemplateDoc, tmplData *TemplateVarData) (*x509.Certificate, uuid.UUID, error) {
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
		Subject:      t.Subject.pkixName(tmplData),
		NotBefore:    now,
		NotAfter:     now.AddDate(0, int(t.ValidityInMonths), 0),
	}
	t.SubjectAlternativeNames.populateCertificate(&c, tmplData)
	if err != nil {
		return nil, certID, err
	}
	if nsType == NSTypeRootCA {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLen = 1
		c.MaxPathLenZero = false
		c.BasicConstraintsValid = true
	} else if nsType == NSTypeIntCA {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLenZero = true
		c.BasicConstraintsValid = true
	} else {
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
		if t.Usage == UsageServerAndClient || t.Usage == UsageServerOnly {
			c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
		if t.Usage == UsageServerAndClient || t.Usage == UsageClientOnly {
			c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
	}

	return &c, certID, err
}

// (nsType/nsID) must be verified prior to calling this function
func (s *adminServer) createCertificateFromTemplate(ctx context.Context, nsType NamespaceTypeShortName, nsID uuid.UUID,
	t *CertificateTemplateDoc, tmplData *TemplateVarData) (*CertDoc, []byte, error) {
	c, certID, err := prepareUnsignedCertificateFromTemplate(nsType, nsID, t, tmplData)
	if err != nil {
		return nil, nil, err
	}

	// prep signer
	signer, err := s.loadCertSigner(ctx, nsType, nsID, t, c, tmplData)
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
		SubjectAlternativeNames: t.SubjectAlternativeNames,
		CertStorePath:           blobName, // certificate storage path in blob storage
		IssuerNamespaceID:       t.IssuerNamespaceID,
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
