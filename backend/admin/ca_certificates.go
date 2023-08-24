package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/stephenzsy/small-kms/backend/common"
)

func (s *adminServer) ListCACertificates(c *gin.Context, params common.ListCACertificatesParams) {
	c.JSON(200, &common.CertificateRefs{})
}

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
	category   common.CreateCertificateParametersCategory
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
	category           common.CreateCertificateParametersCategory
	name               string
	kty                common.CreateCertificateParametersKty
	size               common.CreateCertificateParametersSize
	owner              uuid.UUID
	keyVaultKeyName    string
	keyVaultKeyVersion string
	subject            common.CertificateSubject
}

func getKeyStoreString(keyId *azkeys.ID) string {
	return fmt.Sprintf("%s/%s", keyId.Name(), keyId.Version())
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
}

func (s *adminServer) findLatestCertificate(category common.CreateCertificateParametersCategory, name string) (dto certificateMetadataDto, err error) {
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
	var keyid azkeys.ID
	if len(params.keyVaultKeyVersion) != 0 {
		keyResp, err := keysClient.GetKey(c, params.keyVaultKeyName, params.keyVaultKeyVersion, nil)
		if err != nil {
			log.Error(err)
			return result, err
		}
		keyid = *keyResp.Key.KID
	}
	if len(keyid) == 0 {
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
		keyid = *keyResp.Key.KID
		if err != nil {
			log.Error(err)
			return result, err
		}
	}
	_, err = db.ExecContext(c, `UPDATE cert_metadata
	SET key_store = ?
	WHERE id = ?`, keyid, recordId)
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
