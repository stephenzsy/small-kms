package admin

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type certificateMetadataRow struct {
	id         int
	uuid       string
	category   string
	name       string
	revoked    int
	notBefore  *string
	notAfter   *string
	certStore  *string
	keyStore   *string
	issuer     *int
	owner      *string
	commonName string
}

type certificateMetadataDto struct {
	id         int
	uuid       uuid.UUID
	category   CertificateCategory
	name       string
	revoked    int
	notBefore  time.Time
	notAfter   time.Time
	certStore  string
	keyStore   string
	issuer     int
	owner      string
	commonName string
}

type CertItem struct {
	ID                uuid.UUID           `json:"id"`
	Category          CertificateCategory `json:"category"`
	Name              string              `json:"name"`
	SubjectCommonName string              `json:"subjectCommonName"`
	OwnerID           uuid.UUID           `json:"ownerId"`
	KeyStore          string              `json:"keyStore,omitempty"`
	CertStore         string              `json:"certStore,omitempty"`
	NotAfter          time.Time           `json:"notAfter,omitempty"`
}

func (s *adminServer) ListCertificates(c *gin.Context, category CertificateCategory, params ListCertificatesParams) {
	db := s.config.AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(uuid.Nil.String())
	pager := db.NewQueryItemsPager(`SELECT
		c.id,
		c.category,
		c.name,
		c.subjectCommonName, 
		c.ownerId,
		c.keyStore,
		c.certStore,
		c.notAfter
	 FROM c
	 WHERE c.category = @category AND c.ownerId = @ownerId
	 ORDER BY c.notAfter DESC`, partitionKey, &azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@category", Value: category},
			{Name: "@ownerId", Value: uuid.Nil.String()},
		},
	})
	results := make([]CertificateRef, 0)
	for pager.More() {
		t, err := pager.NextPage(c)
		if err != nil {
			log.Printf("Faild to get list of certificates: %s", err.Error())
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		for _, itemBytes := range t.Items {
			item := CertItem{}
			if err = json.Unmarshal(itemBytes, &item); err != nil {
				log.Printf("Faild to serialize db entry: %s", err.Error())
				c.JSON(500, gin.H{"error": "internal error"})
				return
			}

			results = append(results, CertificateRef{
				CommonName: item.SubjectCommonName,
				ID:         item.ID,
				IssuerID:   item.ID,
				NotAfter:   item.NotAfter,
			})
		}
	}

	c.JSON(200, results)
}

func (dao *certificateMetadataRow) toDTO(dto *certificateMetadataDto) (err error) {
	dto.id = dao.id
	dto.uuid, err = uuid.Parse(dao.uuid)
	if err != nil {
		return
	}
	switch dao.category {
	case string(RootCa):
		dto.category = RootCa
	default:
		err = fmt.Errorf("category not supported: %s", dao.category)
		return
	}
	dto.name = dao.name
	dto.revoked = dao.revoked
	if dao.notBefore != nil && len(*dao.notBefore) > 0 {
		dto.notBefore, err = time.Parse(time.RFC3339, *dao.notBefore)
		if err != nil {
			return
		}
	}
	if dao.notAfter != nil && len(*dao.notAfter) > 0 {
		dto.notAfter, err = time.Parse(time.RFC3339, *dao.notAfter)
		if err != nil {
			return
		}
	}
	if dao.certStore != nil {
		dto.certStore = *dao.certStore
	}
	if dao.keyStore != nil {
		dto.keyStore = *dao.keyStore
	}
	if dao.issuer != nil {
		dto.issuer = *dao.issuer
	}
	if dao.owner != nil {
		dto.owner = *dao.owner
	}
	dto.commonName = dao.commonName
	return
}

type createCertificateInternalParameters struct {
	category           CertificateCategory
	name               string
	kty                CreateCertificateParametersKty
	size               CreateCertificateParametersSize
	owner              uuid.UUID
	keyVaultKeyName    string
	keyVaultKeyVersion string
	subject            CertificateSubject
}

var seededRand *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
}

func (s *adminServer) findLatestCertificate(category CertificateCategory, name string) (dto certificateMetadataDto, err error) {
	return
	db := s.config.GetDB()
	row := db.QueryRow(`SELECT 
		id,
		uuid,
		category,
		name,
		revoked,
		not_before,
		not_after,
		cert_store,
		key_store,
		issuer,
		owner,
		common_name
	FROM cert_metadata WHERE category = ? AND name = ?
	ORDER BY not_after DESC
	LIMIT 1`, category, name)
	dao := certificateMetadataRow{}
	err = row.Scan(
		&dao.id,
		&dao.uuid,
		&dao.category,
		&dao.name,
		&dao.revoked,
		&dao.notBefore,
		&dao.notAfter,
		&dao.certStore,
		&dao.keyStore,
		&dao.issuer,
		&dao.owner,
		&dao.commonName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	err = dao.toDTO(&dto)
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

func (s *adminServer) createCACertificate(c *gin.Context, params createCertificateInternalParameters) (result CertificateRef, err error) {

	// create entry
	db := s.config.AzCosmosContainerClient()
	item := CertItem{
		ID:                uuid.New(),
		Category:          params.category,
		Name:              params.name,
		SubjectCommonName: params.subject.CommonName,
	}
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return
	}
	partitionKey := azcosmos.NewPartitionKeyString(item.OwnerID.String())
	if _, err = db.CreateItem(c, partitionKey, itemBytes, nil); err != nil {
		return
	}

	log.Printf("Created certificate record: %s", item.ID.String())
	// first create new version of key in keyvault
	keysClient := s.config.GetAzKeysClient()
	var webKey *azkeys.JSONWebKey
	if len(params.keyVaultKeyVersion) != 0 {
		keyResp, err := keysClient.GetKey(c, params.keyVaultKeyName, params.keyVaultKeyVersion, nil)
		if err != nil {
			log.Printf("Error getting key: %s", err.Error())
			return result, err
		}
		webKey = keyResp.Key
	}
	if webKey == nil {
		ckp := azkeys.CreateKeyParameters{}
		switch params.kty {
		case RSA:
			ckp.Kty = to.Ptr(azkeys.KeyTypeRSA)

			switch params.size {
			case N4096:
				ckp.KeySize = to.Ptr(int32(4096))
			}
		}
		keyResp, err := keysClient.CreateKey(c, params.keyVaultKeyName, ckp, nil)
		webKey = keyResp.Key

		if err != nil {
			log.Printf("Error getting key: %s", err.Error())
			return result, err
		}
	}

	patchKeyStoreOps := azcosmos.PatchOperations{}
	patchKeyStoreOps.AppendSet("/keyStore", string(*webKey.KID))
	if _, err = db.PatchItem(c, partitionKey, item.ID.String(), patchKeyStoreOps, nil); err != nil {
		return
	} else {
		item.KeyStore = string(*webKey.KID)
	}

	// self-sign

	caSubjectOU := []string{}
	caSubjectO := []string{}
	caSubjectC := []string{}
	if params.subject.OrganizationUnit != nil && len(*params.subject.OrganizationUnit) > 0 {
		caSubjectOU = append(caSubjectOU, *params.subject.OrganizationUnit)
	}
	if params.subject.Organization != nil && len(*params.subject.Organization) > 0 {
		caSubjectO = append(caSubjectO, *params.subject.Organization)
	}
	if params.subject.Country != nil && len(*params.subject.Country) > 0 {
		caSubjectC = append(caSubjectC, *params.subject.Country)
	}
	serialNumber := big.NewInt(0)
	serialNumber = serialNumber.SetBytes(item.ID[:])
	ca := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         params.subject.CommonName,
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
		ctx:        c,
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
	_, err = blobClient.UploadBuffer(c, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", blobKey, "cert.der"), certBytes, nil)
	if err != nil {
		return
	}
	_, err = blobClient.UploadBuffer(c, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", blobKey, "cert.pem"), caPEM.Bytes(), nil)
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
	if _, err = db.PatchItem(c, partitionKey, item.ID.String(), patchCertDoc, nil); err != nil {
		return
	} else {
		item.CertStore = blobKey
		item.NotAfter = parsed.NotAfter.UTC()
	}
	return
}

func (s *adminServer) CreateCertificate(c *gin.Context, params CreateCertificateParams) {
	body := CreateCertificateParameters{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}
	p := createCertificateInternalParameters{
		category: body.Category,
		name:     body.Name,
		subject:  body.Subject,
	}
	lastCertificate, err := s.findLatestCertificate(p.category, p.name)
	if err != nil {
		log.Printf("Error find latest certificate: %s", err.Error())
		c.JSON(500, gin.H{"message": "internal error"})
		return
	}
	if lastCertificate.id > 0 {
		if len(lastCertificate.keyStore) > 0 {
			keyId := azkeys.ID(lastCertificate.keyStore)
			if body.Options != nil && body.Options.KeepKeyVersion != nil && *body.Options.KeepKeyVersion {
				p.keyVaultKeyName = keyId.Name()
				p.keyVaultKeyVersion = keyId.Version()
			} else if body.Options == nil || body.Options.NewKeyName == nil || !*body.Options.NewKeyName {
				p.keyVaultKeyName = keyId.Name()
			}
		}
	}
	switch body.Category {
	case RootCa:
		if body.Kty == nil || len(*body.Kty) == 0 || *body.Kty == RSA {
			p.kty = RSA
			if body.Size == nil || *body.Size == 0 || *body.Size == N4096 {
				p.size = N4096
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
		certCreated, err := s.createCACertificate(c, p)
		if err != nil {
			c.JSON(400, gin.H{"message": "Failed to create certificate", "error": err.Error()})
			log.Printf("Failed to create cert: %s", err.Error())
			return
		}
		c.JSON(201, &certCreated)
	default:
		c.JSON(400, gin.H{"message": "Category not supported", "category": body.Category})
		return
	}
}

func (s *adminServer) DownloadCertificate(c *gin.Context, id uuid.UUID, params DownloadCertificateParams) {
	db := s.config.GetDB()
	row := db.QueryRow(`SELECT 
		cert_store
	FROM cert_metadata WHERE uuid = ?`, id)
	var certStore string
	err := row.Scan(&certStore)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		log.Printf("Faild to get certificates certificate metadata: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if len(certStore) == 0 {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	filename := "cert.pem"
	contentType := "application/x-pem-file"
	if params.Format != nil {
		switch *params.Format {
		case Der:
			filename = "cert.der"
			contentType = "application/x-x509-ca-cert"
		default:
		}
	}

	blobClient := s.config.GetAzBlobClient()

	get, err := blobClient.DownloadStream(c, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", certStore, filename), nil)
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(c, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	err = retryReader.Close()
	if err != nil {
		log.Printf("Faild to get download stream for certificate: %s", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.Data(200, contentType, downloadedData.Bytes())
}
