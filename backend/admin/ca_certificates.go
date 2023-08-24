package admin

import (
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

type createCertificateInternalParameters struct {
	category        common.CreateCertificateParametersCategory
	name            string
	kty             common.CreateCertificateParametersKty
	size            common.CreateCertificateParametersSize
	owner           uuid.UUID
	keyVaultKeyName string
	subject         common.CertificateSubject
}

func getKeyStoreString(keyId *azkeys.ID) string {
	return fmt.Sprintf("%s/%s", keyId.Name(), keyId.Version())
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomHexSuffix(prefix string) string {
	n := seededRand.Int31() % 0x10000
	return fmt.Sprintf("%s%04x", prefix, n)
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
	if err != nil {
		log.Error(err)
		return
	}
	_, err = db.ExecContext(c, `UPDATE cert_metadata
		SET key_store = ?
		WHERE id = ?`, getKeyStoreString(keyResp.Key.KID), recordId)
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
		p.keyVaultKeyName = generateRandomHexSuffix("root-ca-")
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
