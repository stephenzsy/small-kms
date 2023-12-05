package cert

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	pkcs12utils "github.com/stephenzsy/small-kms/backend/internal/pkcs12"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// ExchangePKCS12 implements admin.ServerInterface.
func (*CertServer) ExchangePKCS12(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, id string) error {
	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)

	if _, authOk := authz.Authorize(c, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	req := new(certmodels.ExchangePKCS12Request)
	if err := c.Bind(req); err != nil {
		return err
	}

	jwe, err := cloudkey.NewJsonWebEncryption(req.Payload)
	if err != nil {
		return fmt.Errorf("%w: invalid payload", base.ErrResponseStatusBadRequest)
	}
	certDoc := certDocInternal{}
	if err := readCertDocInternal(c, namespaceProvider, namespaceId, id, &certDoc); err != nil {
		return err
	}
	otk := certDoc.ExportKey
	if otk == nil {
		return base.ErrResponseStatusBadRequest
	}
	deleteOtk := false
	defer func() {
		if deleteOtk {
			ec := c.Elevate()
			docSvc := resdoc.GetDocService(ec)
			patchOps := azcosmos.PatchOperations{}
			patchOps.AppendRemove("/exportKey")
			_, err := docSvc.Patch(ec, &certDoc, patchOps, nil)
			if err != nil {
				log.Ctx(ec).Err(err).Msg("Error delete export key")
			}
		}
	}()

	reqPayload, encKey, err := jwe.Decrypt(func(*cloudkey.JoseHeader) crypto.PrivateKey {
		key, _ := otk.PrivateKey().(*ecdsa.PrivateKey).ECDH()
		return key
	})
	if err != nil {
		return err
	}
	privateJwk := new(cloudkey.JsonWebKey)
	if err := json.Unmarshal(reqPayload, privateJwk); err != nil {
		return fmt.Errorf("%w: invalid payload", base.ErrResponseStatusBadRequest)
	}

	certChain := make([]*x509.Certificate, len(certDoc.JsonWebKey.CertificateChain))
	for i, certBytes := range certDoc.JsonWebKey.CertificateChain {
		if certChain[i], err = x509.ParseCertificate(certBytes); err != nil {
			return fmt.Errorf("%w: invalid certificate chain", base.ErrResponseStatusBadRequest)
		}
	}
	legacy := false
	if req.Legacy != nil && *req.Legacy {
		legacy = true
	}

	password := ""
	if req.PasswordProtected {
		randUUID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		password = randUUID.String()[0:8]
	}

	pkcs12File, err := pkcs12utils.ConvertPKCS12(privateJwk.PrivateKey(), certChain, password, legacy)
	if err != nil {
		return err
	}

	var resultPayload string

	resultHeader := jwe.Protected
	switch jwe.Protected.Algorithm {
	case cloudkey.JwkEncAlgEcdhEs:
		resultHeader = cloudkey.JoseHeader{
			EncryptionAlgorithm: jwe.Protected.EncryptionAlgorithm,
		}
		if headerJson, err := json.Marshal(resultHeader); err != nil {
			return err
		} else {
			resultHeader.Raw = base64.RawURLEncoding.EncodeToString(headerJson)
		}
	}

	switch jwe.Protected.EncryptionAlgorithm {
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
		encrypted := gcm.Seal(nil, iv, pkcs12File, []byte(resultHeader.Raw))
		ciphertext := encrypted[:len(encrypted)-ci.BlockSize()]
		tag := encrypted[len(encrypted)-ci.BlockSize():]
		resultJwe := &cloudkey.JsonWebEncryption{
			Protected:            resultHeader,
			EncryptedKey:         jwe.EncryptedKey,
			InitializationVector: iv,
			Ciphertext:           ciphertext,
			AuthenticationTag:    tag,
		}
		resultPayload = resultJwe.String()
	}

	deleteOtk = true
	return c.JSON(http.StatusOK, &certmodels.ExchangePKCS12Result{
		Payload:  resultPayload,
		Password: password,
	})
}
