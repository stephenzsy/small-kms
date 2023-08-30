package admin

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	mrand "math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/google/uuid"
)

var seededRand *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
}

func toPublicRSA(key *azkeys.JSONWebKey) (*rsa.PublicKey, error) {
	res := &rsa.PublicKey{}

	// N = modulus
	if len(key.N) == 0 {
		return nil, errors.New("property N is empty")
	}
	res.N = &big.Int{}
	res.N.SetBytes(key.N)

	// e = public exponent
	if len(key.E) == 0 {
		return nil, errors.New("property e is empty")
	}
	res.E = int(big.NewInt(0).SetBytes(key.E).Uint64())
	return res, nil
}

func (s *adminServer) createCACertificate(ctx context.Context, p createCertificateInternalParameters) (item CertDBItem, err error) {

	// create entry
	db := s.azCosmosContainerClientCerts
	certId := uuid.New()
	item.CertificateRef = CertificateRef{
		ID:              certId,
		IssuerNamespace: wellKnownNamespaceID_RootCA,
		Issuer:          p.issuer.ID,
		Usage:           p.usage,
		Name:            p.subject.CN,
		NamespaceID:     p.namespaceID,
		CreatedBy:       p.createdBy,
	}
	if p.usage == UsageRootCA {
		item.CertificateRef.Issuer = certId
	}

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return
	}
	partitionKey := azcosmos.NewPartitionKeyString(p.namespaceID.String())
	if _, err = db.CreateItem(ctx, partitionKey, itemBytes, nil); err != nil {
		return
	}
	log.Printf("Created certificate record: %s", item.ID.String())

	// first create new version of key in keyvault
	var webKey *azkeys.JSONWebKey
	if len(p.keyVaultKeyVersion) != 0 {
		keyResp, err := s.azKeysClient.GetKey(ctx, p.keyVaultKeyName, p.keyVaultKeyVersion, nil)
		if err != nil {
			log.Printf("Error getting key: %s", err.Error())
			return item, err
		}
		webKey = keyResp.Key
	}
	if webKey == nil {
		ckp := azkeys.CreateKeyParameters{}
		switch p.kty {
		case KtyRSA:
			ckp.Kty = to.Ptr(azkeys.KeyTypeRSA)

			switch p.size {
			case KeySize4096:
				ckp.KeySize = to.Ptr(int32(4096))
			case KeySize2048:
				ckp.KeySize = to.Ptr(int32(2048))
			}
		}
		keyResp, err := s.azKeysClient.CreateKey(ctx, p.keyVaultKeyName, ckp, nil)
		webKey = keyResp.Key

		if err != nil {
			log.Printf("Error getting key: %s", err.Error())
			return item, err
		}
	}

	patchKeyStoreOps := azcosmos.PatchOperations{}
	patchKeyStoreOps.AppendSet("/keyStore", string(*webKey.KID))
	if _, err = db.PatchItem(ctx, partitionKey, item.ID.String(), patchKeyStoreOps, nil); err != nil {
		return
	} else {
		item.KeyStore = string(*webKey.KID)
	}

	// self-sign

	caSubjectOU := []string{}
	caSubjectO := []string{}
	caSubjectC := []string{}
	if p.subject.OU != nil && len(*p.subject.OU) > 0 {
		caSubjectOU = append(caSubjectOU, *p.subject.OU)
	}
	if p.subject.O != nil && len(*p.subject.O) > 0 {
		caSubjectO = append(caSubjectO, *p.subject.O)
	}
	if p.subject.C != nil && len(*p.subject.C) > 0 {
		caSubjectC = append(caSubjectC, *p.subject.C)
	}
	serialNumber := big.NewInt(0)
	serialNumber = serialNumber.SetBytes(item.ID[:])

	ca := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         p.subject.CN,
			OrganizationalUnit: caSubjectOU,
			Organization:       caSubjectO,
			Country:            caSubjectC,
		},
		NotBefore:             time.Now(),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		MaxPathLenZero:        true,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA384WithRSA,
	}

	var issuerCertPemBlock *pem.Block = nil
	var signerCert *x509.Certificate = nil
	if p.usage == UsageRootCA {
		ca.MaxPathLen = 1
		ca.MaxPathLenZero = false
		ca.NotAfter = ca.NotBefore.AddDate(10, 0, 0)
		signerCert = &ca
	} else {
		// load signer
		pemBlob, err1 := s.FetchCertificatePEMBlob(ctx, p.issuer.CertStore)
		if err1 != nil {
			err = err1
			return
		}
		issuerCertPemBlock, _ = pem.Decode(pemBlob)
		if issuerCertPemBlock == nil {
			err = errors.New("failed to decode issuer certificate")
			return
		}
		if signerCert, err = x509.ParseCertificate(issuerCertPemBlock.Bytes); err != nil {
			return
		}
		if signerCert == nil {
			err = fmt.Errorf("unable to load issuer certificate: %s/%s", p.issuer.IssuerNamespace.String(), p.issuer.ID.String())
			return
		}

		if signerCert.NotAfter.Before(ca.NotBefore.AddDate(0, 0, 15)) {
			err = fmt.Errorf("issuer certificate is about to expire in 15 days: %s/%s", p.issuer.IssuerNamespace.String(), p.issuer.ID.String())
			return
		}
		if p.namespaceID == wellKnownNamespaceID_IntCAService {
			ca.NotAfter = ca.NotBefore.AddDate(5, 0, 0)
		} else if p.namespaceID == wellKnownNamespaceID_IntCaSCEPIntranet {
			ca.NotAfter = ca.NotBefore.AddDate(1, 0, 0)
		}
		if ca.NotAfter.After(signerCert.NotAfter) {
			ca.NotAfter = signerCert.NotAfter
		}
	}

	// public key
	pubKey, err := toPublicRSA(webKey)
	if err != nil {
		return
	}

	var signerPubKey crypto.PublicKey
	if p.usage == UsageRootCA {
		signerPubKey = pubKey
	} else {
		signerPubKey = signerCert.PublicKey
	}

	var signerKID azkeys.ID
	if p.usage == UsageRootCA {
		signerKID = *webKey.KID
	} else {
		signerKID = azkeys.ID(p.issuer.KeyStore)
	}
	signer := keyVaultSigner{
		ctx:        ctx,
		keysClient: s.azKeysClient,
		kid:        signerKID,
		publicKey:  signerPubKey,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &ca, signerCert, pubKey, &signer)
	if err != nil {
		return
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return
	}
	if issuerCertPemBlock != nil {
		err = pem.Encode(caPEM, issuerCertPemBlock)
		if err != nil {
			return
		}
	}

	blobName := fmt.Sprintf("%s/%s.pem", p.keyVaultKeyName, item.ID)
	// upload to blob storage
	blobClient := s.azBlobContainerClient.NewBlockBlobClient(blobName)
	_, err = blobClient.UploadBuffer(ctx, caPEM.Bytes(), &blockblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr("application/x-pem-file"),
		},
	})
	if err != nil {
		return
	}

	parsed, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	patchCertDoc := azcosmos.PatchOperations{}
	patchCertDoc.AppendSet("/certStore", blobClient.URL())
	patchCertDoc.AppendSet("/notAfter", parsed.NotAfter.UTC().Format(time.RFC3339))
	if _, err = db.PatchItem(ctx, partitionKey, item.ID.String(), patchCertDoc, nil); err != nil {
		return
	} else {
		item.CertStore = blobName
		item.NotAfter = parsed.NotAfter.UTC()
	}

	return
}
