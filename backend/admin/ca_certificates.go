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
	"io"
	"log"
	"math/big"
	mrand "math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
)

type createCertificateInternalParameters struct {
	usage              CertificateUsage
	kty                CreateCertificateParametersKty
	size               CreateCertificateParametersSize
	namespaceID        uuid.UUID
	keyVaultKeyName    string
	keyVaultKeyVersion string
	subject            CertificateSubject
	createdBy          string
}

var seededRand *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
}

func (s *adminServer) findLatestCertificate(ctx context.Context, usage CertificateUsage, name string) (result CertDBItem, err error) {
	partitionKey := azcosmos.NewPartitionKeyString(string(WellKnownNamespaceIDStrRootCA))
	db := s.config.AzCosmosContainerClient()
	pager := db.NewQueryItemsPager(`
SELECT TOP 1
	*
FROM c
WHERE c.usage = @usage AND c.name = @name
ORDER BY c.notAfter DESC`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@usage", Value: usage},
				{Name: "@name", Value: name},
			},
		})
	t, err := pager.NextPage(ctx)
	if err != nil {
		return
	}
	if len(t.Items) > 0 {
		err = json.Unmarshal(t.Items[0], &result)
	}
	return
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

type keyVaultSigner struct {
	ctx        context.Context
	keysClient *azkeys.Client
	webKey     *azkeys.JSONWebKey
	publicKey  crypto.PublicKey
}

func (s *keyVaultSigner) Public() crypto.PublicKey {
	return s.publicKey
}

func (s *keyVaultSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	resp, err := s.keysClient.Sign(s.ctx, s.webKey.KID.Name(), s.webKey.KID.Version(), azkeys.SignParameters{
		Algorithm: to.Ptr(azkeys.SignatureAlgorithmRS384),
		Value:     digest,
	}, nil)
	if err != nil {
		return
	}
	signature = resp.Result
	return
}

func (s *adminServer) createRootCACertificate(ctx context.Context, params createCertificateInternalParameters) (item CertDBItem, err error) {

	// create entry
	db := s.config.AzCosmosContainerClient()
	certId := uuid.New()
	item.CertificateRef = CertificateRef{
		ID:              certId,
		Issuer:          certId,
		IssuerNamespace: wellKnownNamespaceIDRootCA,
		Usage:           UsageRootCA,
		Name:            params.subject.CN,
		NamespaceID:     params.namespaceID,
		CreatedBy:       params.createdBy,
	}
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return
	}
	partitionKey := azcosmos.NewPartitionKeyString(wellKnownNamespaceIDRootCA.String())
	if _, err = db.CreateItem(ctx, partitionKey, itemBytes, nil); err != nil {
		return
	}
	log.Printf("Created certificate record: %s", item.ID.String())

	// first create new version of key in keyvault
	keysClient := s.config.GetAzKeysClient()
	var webKey *azkeys.JSONWebKey
	if len(params.keyVaultKeyVersion) != 0 {
		keyResp, err := keysClient.GetKey(ctx, params.keyVaultKeyName, params.keyVaultKeyVersion, nil)
		if err != nil {
			log.Printf("Error getting key: %s", err.Error())
			return item, err
		}
		webKey = keyResp.Key
	}
	if webKey == nil {
		ckp := azkeys.CreateKeyParameters{}
		switch params.kty {
		case KtyRSA:
			ckp.Kty = to.Ptr(azkeys.KeyTypeRSA)

			switch params.size {
			case KeySize4096:
				ckp.KeySize = to.Ptr(int32(4096))
			}
		}
		keyResp, err := keysClient.CreateKey(ctx, params.keyVaultKeyName, ckp, nil)
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
	if params.subject.OU != nil && len(*params.subject.OU) > 0 {
		caSubjectOU = append(caSubjectOU, *params.subject.OU)
	}
	if params.subject.O != nil && len(*params.subject.O) > 0 {
		caSubjectO = append(caSubjectO, *params.subject.O)
	}
	if params.subject.C != nil && len(*params.subject.C) > 0 {
		caSubjectC = append(caSubjectC, *params.subject.C)
	}
	serialNumber := big.NewInt(0)
	serialNumber = serialNumber.SetBytes(item.ID[:])
	ca := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         params.subject.CN,
			OrganizationalUnit: caSubjectOU,
			Organization:       caSubjectO,
			Country:            caSubjectC,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		MaxPathLen:            1,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA384WithRSA,
	}

	// public key
	pubKey, err := toPublicRSA(webKey)
	if err != nil {
		return
	}

	signer := keyVaultSigner{
		ctx:        ctx,
		keysClient: keysClient,
		webKey:     webKey,
		publicKey:  pubKey,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &ca, &ca, pubKey, &signer)
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

	blobKey := fmt.Sprintf("%s/%s", params.keyVaultKeyName, item.ID)
	// upload to blob storage
	blobClient := s.config.GetAzBlobClient()
	_, err = blobClient.UploadBuffer(ctx, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", blobKey, "cert.der"), certBytes, nil)
	if err != nil {
		return
	}
	_, err = blobClient.UploadBuffer(ctx, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", blobKey, "cert.pem"), caPEM.Bytes(), nil)
	if err != nil {
		return
	}

	parsed, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	patchCertDoc := azcosmos.PatchOperations{}
	patchCertDoc.AppendSet("/certStore", blobKey)
	patchCertDoc.AppendSet("/notAfter", parsed.NotAfter.UTC().Format(time.RFC3339))
	if _, err = db.PatchItem(ctx, partitionKey, item.ID.String(), patchCertDoc, nil); err != nil {
		return
	} else {
		item.CertStore = blobKey
		item.NotAfter = parsed.NotAfter.UTC()
	}
	return
}

func (s *adminServer) CreateCertificateV1(c *gin.Context, namespaceID NamespaceID) {
	body := CreateCertificateParameters{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}
	if namespaceID != wellKnownNamespaceIDRootCA {
		c.JSON(400, gin.H{"message": "invalid namespace", "namespace": namespaceID})
		return
	}

	if !auth.HasAdminAppRole(c) {
		c.JSON(403, gin.H{"message": "User must have admin role"})
		return
	}

	p := createCertificateInternalParameters{
		usage:       body.Usage,
		subject:     body.Subject,
		namespaceID: namespaceID,
		createdBy:   auth.GetCallerID(c),
	}
	lastCertificate, err := s.findLatestCertificate(c.Request.Context(), p.usage, p.subject.CN)
	if err != nil {
		log.Printf("Error find latest certificate: %s", err.Error())
		c.JSON(500, gin.H{"message": "internal error"})
		return
	}
	if lastCertificate.ID != uuid.Nil {
		if len(lastCertificate.KeyStore) > 0 {
			keyId := azkeys.ID(lastCertificate.KeyStore)
			if body.Options != nil && body.Options.KeepKeyVersion != nil && *body.Options.KeepKeyVersion {
				p.keyVaultKeyName = keyId.Name()
				p.keyVaultKeyVersion = keyId.Version()
			} else if body.Options == nil || body.Options.NewKeyName == nil || !*body.Options.NewKeyName {
				p.keyVaultKeyName = keyId.Name()
			}
		}
	}
	switch body.Usage {
	case body.Usage:
		if body.Kty == nil || len(*body.Kty) == 0 || *body.Kty == KtyRSA {
			p.kty = KtyRSA
			if body.Size == nil || *body.Size == 0 || *body.Size == KeySize4096 {
				p.size = KeySize4096
			} else {
				c.JSON(400, gin.H{"message": "Size not supported", "size": body.Size})
				return
			}
		} else {
			c.JSON(400, gin.H{"message": "Key type not supported", "kty": body.Kty})
			return
		}
		if len(p.keyVaultKeyName) == 0 {
			p.keyVaultKeyName = generateRandomHexSuffix("root-ca-")
		}
		certCreated, err := s.createRootCACertificate(c, p)
		if err != nil {
			c.JSON(400, gin.H{"message": "Failed to create certificate", "error": err.Error()})
			log.Printf("Failed to create cert: %s", err.Error())
			return
		}
		c.JSON(201, &certCreated.CertificateRef)
	default:
		c.JSON(400, gin.H{"message": "Usage not supported", "usage": body.Usage})
		return
	}
}
