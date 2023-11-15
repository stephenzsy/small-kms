package managedapp

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/admin/systemapp"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	pkcs12utils "github.com/stephenzsy/small-kms/backend/internal/pkcs12"
	"github.com/stephenzsy/small-kms/backend/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// ExchangePKCS12 implements ServerInterface.
func (s *server) ExchangePKCS12(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, certID base.ID) error {
	c := ec.(ctx.RequestContext)
	logger := log.Ctx(c)
	c, nsCtx := ns.WithResovingMeNSContext(c, namespaceKind, namespaceId)

	c, authOK := authz.Authorize(c, nsCtx.AllowSelf())
	if !authOK {
		return base.ErrResponseStatusForbidden
	}

	req := new(ExchangePKCS12Request)
	if err := c.Bind(req); err != nil {
		return err
	}

	jwe, err := cloudkey.NewJsonWebEncryption(req.Payload)
	if err != nil {
		return fmt.Errorf("%w: invalid payload", base.ErrResponseStatusBadRequest)
	}
	switch jwe.Protected.Enc {
	case cloudkey.JwkEncAlgAes256Gcm:
	default:
		return fmt.Errorf("%w: unsupported algorithm: %s", base.ErrResponseStatusBadRequest, jwe.Protected.Enc)
	}

	certDoc, err := cert.ApiReadCertDocByID(c, certID)
	if err != nil {
		return err
	}

	if time.Now().After(certDoc.Created.Time.Add(time.Hour)) {
		return fmt.Errorf("%w: this operation can only be performed within 1 hour after cert issurance", base.ErrResponseStatusBadRequest)
	}

	systemAppDoc, _, err := systemapp.GetSystemAppDoc(c, systemapp.SystemAppNameBackend)
	if err != nil {
		return err
	}
	if systemAppDoc.ServicePrincipalID == "" ||
		req.KeyLocator.NamespaceKind() != base.NamespaceKindServicePrincipal ||
		req.KeyLocator.NamespaceID() != base.ParseID(systemAppDoc.ServicePrincipalID) {
		return fmt.Errorf("%w: invalid key locator", base.ErrResponseStatusBadRequest)
	}

	keyDoc, err := key.ApiReadKeyDocByLocator(c, *req.KeyLocator)
	if err != nil {
		return err
	}

	if keyDoc.KeyID != jwe.Protected.KeyID {
		return fmt.Errorf("%w: keyIDs must match, unwrap keyID: %s, provided: %s", base.ErrResponseStatusBadRequest, keyDoc.KeyID, jwe.Protected.KeyID)
	}

	cloudKey := cloudkeyaz.NewCloudWrappingKeyWithKID(c, kv.GetAzKeyVaultService(c).AzKeysClient(), keyDoc.KeyID, keyDoc.KeyType)
	decryptedPayload, encKey, err := jwe.Decrypt(func(_ *cloudkey.JoseHeader) crypto.Decrypter {
		return cloudKey
	})
	if err != nil {
		logger.Warn().Err(err).Msg("failed to decrypt payload")
		return fmt.Errorf("%w: failed to decrypt payload", base.ErrResponseStatusBadRequest)
	}
	payload := new(RequestPayload)
	if err != json.Unmarshal(decryptedPayload, payload) {
		return fmt.Errorf("%w: invalid payload", base.ErrResponseStatusBadRequest)
	}

	certChain := make([]*x509.Certificate, len(certDoc.KeySpec.CertificateChain))
	for i, certBytes := range certDoc.KeySpec.CertificateChain {
		if certChain[i], err = x509.ParseCertificate(certBytes); err != nil {
			return fmt.Errorf("%w: invalid certificate chain", base.ErrResponseStatusBadRequest)
		}
	}
	legacy := false
	if req.Legacy != nil && *req.Legacy {
		legacy = true
	}
	pkcs12File, err := pkcs12utils.ConvertPKCS12(payload.PrivateKey.PrivateKey(), certChain, payload.Password, legacy)
	if err != nil {
		return err
	}

	var resultPayload string

	switch jwe.Protected.Enc {
	case cloudkey.JwkEncAlgAes256Gcm:
		ci, err := aes.NewCipher(encKey)
		if err != nil {
			return err
		}
		gcm, err := cipher.NewGCM(ci)
		if err != nil {
			return err
		}
		iv := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(iv); err != nil {
			return err
		}
		encrypted := gcm.Seal(nil, iv, pkcs12File, []byte(jwe.Protected.Raw))
		ciphertext := encrypted[:len(encrypted)-ci.BlockSize()]
		tag := encrypted[len(encrypted)-ci.BlockSize():]
		resultJwe := &cloudkey.JsonWebEncryption{
			Protected:            jwe.Protected,
			EncryptedKey:         jwe.EncryptedKey,
			InitializationVector: iv,
			Ciphertext:           ciphertext,
			AuthenticationTag:    tag,
		}
		resultPayload = resultJwe.String()
	}

	return c.JSON(http.StatusOK, &ExchangePKCS12Result{
		Payload: resultPayload,
	})
}

type RequestPayload struct {
	PrivateKey cloudkey.JsonWebSignatureKey `json:"privateKey"`
	Password   string                       `json:"password"`
}
