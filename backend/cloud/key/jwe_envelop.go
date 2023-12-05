package cloudkey

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type JWEAes256GcmEncBuilder struct {
	JsonWebEncryption
	encKey []byte
}

func (b *JWEAes256GcmEncBuilder) SetEcdhEsKeyAgreement(selfJWK *JsonWebKey, remoteJWK *JsonWebKey) error {
	b.Protected.EphemeralPublicKey = selfJWK.PublicJWK()
	b.Protected.Algorithm = JwkEncAlgEcdhEs
	b.Protected.EncryptionAlgorithm = JwkEncAlgAes256Gcm

	selfEcdhKey, err := selfJWK.PrivateKey().(*ecdsa.PrivateKey).ECDH()
	if err != nil {
		return err
	}
	remotePublicKey, err := remoteJWK.PublicKey().(*ecdsa.PublicKey).ECDH()
	if err != nil {
		return err
	}

	z, err := selfEcdhKey.ECDH(remotePublicKey)
	if err != nil {
		return err
	}
	kdf := &ecdhesKDF{
		z:   z[:32],
		alg: string(b.Protected.EncryptionAlgorithm),
		apu: b.Protected.AgreementPartyUInfo,
		apv: b.Protected.AgreementPartyVInfo,
	}
	b.encKey = kdf.getAESGCM256DerivedKey()
	return nil
}

func (b *JWEAes256GcmEncBuilder) SetDirectEncryptionKey(key []byte) {
	b.Protected.Algorithm = JwkEncAlgDir
	b.Protected.EncryptionAlgorithm = JwkEncAlgAes256Gcm
	b.encKey = key
}

func (b *JWEAes256GcmEncBuilder) Seal(plaintext []byte) (string, error) {
	if len(b.encKey) == 0 {
		return "", fmt.Errorf("encryption key not set")
	}
	ci, err := aes.NewCipher(b.encKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(ci)
	if err != nil {
		return "", err
	}
	iv := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	if headerJson, err := json.Marshal(b.Protected); err != nil {
		return "", err
	} else {
		b.Protected.Raw = base64.RawURLEncoding.EncodeToString(headerJson)
	}
	encrypted := gcm.Seal(nil, iv, plaintext, []byte(b.Protected.Raw))
	ciphertext := encrypted[:len(encrypted)-ci.BlockSize()]
	tag := encrypted[len(encrypted)-ci.BlockSize():]
	b.InitializationVector = iv
	b.Ciphertext = ciphertext
	b.AuthenticationTag = tag
	return b.String(), nil
}
