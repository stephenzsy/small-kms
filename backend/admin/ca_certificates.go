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
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/stephenzsy/small-kms/backend/common"
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
	category   common.CertificateCategory
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

func (s *adminServer) ListCertificates(c *gin.Context, category common.CertificateCategory, params common.ListCertificatesParams) {
	db := s.config.GetDB()
	rows, err := db.QueryContext(c, `SELECT 
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
		FROM cert_metadata WHERE category = ? AND revoked = 0
		ORDER BY not_after DESC`, category)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(200, &common.CertificateRefs{})
			return
		}
		log.Errorf("Faild to get list of certificates: %w", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	l := make([]common.CertificateRef, 0)
	for rows.Next() {
		dao := certificateMetadataRow{}
		err = rows.Scan(
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
			log.Errorf("Faild to get list of certificates: %w", err)
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		dto := certificateMetadataDto{}
		err = dao.toDTO(&dto)
		if err != nil {
			log.Errorf("Faild to get list of certificates: %w", err)
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}
		l = append(l, common.CertificateRef{
			SerialNumber: fmt.Sprintf("%d", dto.id),
			ID:           dto.uuid,
			IssuerID:     dto.uuid,
			NotAfter:     dto.notAfter,
			NotBefore:    dto.notBefore,
			CommonName:   dto.commonName,
		})
	}
	c.JSON(200, l)
}

func (dao *certificateMetadataRow) toDTO(dto *certificateMetadataDto) (err error) {
	dto.id = dao.id
	dto.uuid, err = uuid.Parse(dao.uuid)
	if err != nil {
		return
	}
	switch dao.category {
	case string(common.RootCa):
		dto.category = common.RootCa
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
	category           common.CertificateCategory
	name               string
	kty                common.CreateCertificateParametersKty
	size               common.CreateCertificateParametersSize
	owner              uuid.UUID
	keyVaultKeyName    string
	keyVaultKeyVersion string
	subject            common.CertificateSubject
}

var seededRand *mrand.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
}

func (s *adminServer) findLatestCertificate(category common.CertificateCategory, name string) (dto certificateMetadataDto, err error) {
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

func (s *adminServer) createCACertificate(c *gin.Context, params createCertificateInternalParameters) (result common.CertificateRef, err error) {
	certUuid := uuid.New()
	// create entry

	db := s.config.GetDB()
	sqlResult, err := db.ExecContext(c, `INSERT INTO cert_metadata(
		uuid, category, name, common_name) VALUES (?, ?, ?, ?)`, certUuid, params.category, params.name, params.subject.CommonName)
	if err != nil {
		return
	}

	recordId, err := sqlResult.LastInsertId()
	if err != nil {
		return
	}
	log.Debugf("Created certificate record %d", recordId)

	// first create new version of key in keyvault
	keysClient := s.config.GetAzKeysClient()
	var webKey *azkeys.JSONWebKey
	if len(params.keyVaultKeyVersion) != 0 {
		keyResp, err := keysClient.GetKey(c, params.keyVaultKeyName, params.keyVaultKeyVersion, nil)
		if err != nil {
			log.Error(err)
			return result, err
		}
		webKey = keyResp.Key
	}
	if webKey == nil {
		ckp := azkeys.CreateKeyParameters{}
		switch params.kty {
		case common.RSA:
			ckp.Kty = to.Ptr(azkeys.KeyTypeRSA)

			switch params.size {
			case common.N4096:
				ckp.KeySize = to.Ptr(int32(4096))
			}
		}
		keyResp, err := keysClient.CreateKey(c, params.keyVaultKeyName, ckp, nil)
		webKey = keyResp.Key

		if err != nil {
			log.Error(err)
			return result, err
		}
	}
	_, err = db.ExecContext(c, `UPDATE cert_metadata
	SET key_store = ?
	WHERE id = ?`, webKey.KID, recordId)
	if err != nil {
		return
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
	ca := x509.Certificate{
		SerialNumber: big.NewInt(recordId),
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

	blobKey := fmt.Sprintf("%s/%d", params.keyVaultKeyName, recordId)
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

	// update record
	_, err = db.ExecContext(c, `UPDATE cert_metadata
	SET cert_store = ?,
		not_before = ?,
		not_after = ?
	WHERE id = ?`, blobKey, parsed.NotBefore.UTC().Format(time.RFC3339), parsed.NotAfter.UTC().Format(time.RFC3339), recordId)
	if err != nil {
		return
	}
	return
}

func (s *adminServer) CreateCertificate(c *gin.Context, params common.CreateCertificateParams) {
	body := common.CreateCertificateParameters{}
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
		log.Errorf("Error find latest certificate: %w", err)
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
	case common.RootCa:
		if body.Kty == nil || len(*body.Kty) == 0 || *body.Kty == common.RSA {
			p.kty = common.RSA
			if body.Size == nil || *body.Size == 0 || *body.Size == common.N4096 {
				p.size = common.N4096
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
			log.Errorf("Failed to create cert: %w", err)
			return
		}
		c.JSON(201, &certCreated)
	default:
		c.JSON(400, gin.H{"message": "Category not supported", "category": body.Category})
		return
	}
}

func (s *adminServer) DownloadCertificate(c *gin.Context, id uuid.UUID, params common.DownloadCertificateParams) {
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
		log.Errorf("Faild to get certificates certificate metadata: %w", err)
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
		case common.Der:
			filename = "cert.der"
			contentType = "application/x-x509-ca-cert"
		default:
		}
	}

	blobClient := s.config.GetAzBlobClient()

	get, err := blobClient.DownloadStream(c, s.config.GetAzBlobContainerName(), fmt.Sprintf("%s/%s", certStore, filename), nil)
	if err != nil {
		log.Errorf("Faild to get download stream for certificate: %w", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(c, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		log.Errorf("Faild to get download stream for certificate: %w", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	err = retryReader.Close()
	if err != nil {
		log.Errorf("Faild to get download stream for certificate: %w", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.Data(200, contentType, downloadedData.Bytes())
}
