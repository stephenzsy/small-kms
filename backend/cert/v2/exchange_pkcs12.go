package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	pkcs12utils "github.com/stephenzsy/small-kms/backend/internal/pkcs12"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
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

	if jwe.Protected.KeyID == "" {
		return fmt.Errorf("%w: invalid payload, one time key id must be specified", base.ErrResponseStatusBadRequest)
	}
	otk, err := key.ReadOneTimeKey(c, namespaceProvider, namespaceId, jwe.Protected.KeyID)
	if err != nil {
		return err
	}

	certDoc := certDocInternal{}
	if err := readCertDocInternal(c, namespaceProvider, namespaceId, id, &certDoc); err != nil {
		return err
	}

	reqPayload, encKey, err := jwe.Decrypt(func(*cloudkey.JoseHeader) (crypto.PrivateKey, error) {
		return otk.PrivateKey().(*ecdsa.PrivateKey).ECDH()
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

	jweBuilder := &cloudkey.JWEAes256GcmEncBuilder{}
	jweBuilder.SetDirectEncryptionKey(encKey)

	resultPayload, err = jweBuilder.Seal(pkcs12File)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &certmodels.ExchangePKCS12Result{
		Payload:  resultPayload,
		Password: password,
	})
}
