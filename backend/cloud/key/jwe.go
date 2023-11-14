package cloudkey

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidJWE = errors.New("invalid JWE")
)

type JsonWebKeyEncryptionAlgorithm string

const (
	JwkEncAlgRsaOeap256 JsonWebKeyEncryptionAlgorithm = "RSA-OAEP-256"
	JwkEncAlgAes256Gcm  JsonWebKeyEncryptionAlgorithm = "A256GCM"
)

type JoseHeader struct {
	Alg   JsonWebKeyEncryptionAlgorithm `json:"alg"`
	Enc   JsonWebKeyEncryptionAlgorithm `json:"enc"`
	KeyID string                        `json:"kid"`
	Raw   string                        `json:"-"`
}

type JsonWebEncryption struct {
	Protected            JoseHeader
	EncryptedKey         Base64RawURLEncodableBytes
	InitializationVector Base64RawURLEncodableBytes
	Ciphertext           Base64RawURLEncodableBytes
	AuthenticationTag    Base64RawURLEncodableBytes
}

func NewJsonWebEncryption(text string) (*JsonWebEncryption, error) {
	parts := strings.Split(text, ".")
	if len(parts) != 5 {
		return nil, ErrInvalidJWE
	}
	jwe := &JsonWebEncryption{
		EncryptedKey:         Base64RawURLEncodableBytes(parts[1]),
		InitializationVector: Base64RawURLEncodableBytes(parts[2]),
		Ciphertext:           Base64RawURLEncodableBytes(parts[3]),
		AuthenticationTag:    Base64RawURLEncodableBytes(parts[4]),
	}
	if decoded, err := base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return nil, err
	} else if err := json.Unmarshal(decoded, &jwe.Protected); err != nil {
		return nil, err
	} else {
		jwe.Protected.Raw = parts[0]
	}

	var err error
	if jwe.EncryptedKey, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return nil, err
	}
	if jwe.InitializationVector, err = base64.RawURLEncoding.DecodeString(parts[2]); err != nil {
		return nil, err
	}
	if jwe.Ciphertext, err = base64.RawURLEncoding.DecodeString(parts[3]); err != nil {
		return nil, err
	}
	if jwe.AuthenticationTag, err = base64.RawURLEncoding.DecodeString(parts[4]); err != nil {
		return nil, err
	}
	return jwe, nil
}

func (jwe *JsonWebEncryption) String() string {
	headerJson, _ := json.Marshal(jwe.Protected)
	sb := strings.Builder{}
	sb.WriteString(base64.RawURLEncoding.EncodeToString(headerJson))
	if text, err := jwe.EncryptedKey.MarshalText(); err == nil {
		sb.WriteString(".")
		sb.Write(text)
	}
	if text, err := jwe.InitializationVector.MarshalText(); err == nil {
		sb.WriteString(".")
		sb.Write(text)
	}
	if text, err := jwe.Ciphertext.MarshalText(); err == nil {
		sb.WriteString(".")
		sb.Write(text)
	}
	if text, err := jwe.AuthenticationTag.MarshalText(); err == nil {
		sb.WriteString(".")
		sb.Write(text)
	}
	return sb.String()
}

func (jwe *JsonWebEncryption) Decrypt(keyFunc func(header *JoseHeader) crypto.Decrypter) (plaintext []byte, unwrappedKey []byte, err error) {

	unwrappedKey, err = jwe.unwrapKey(keyFunc)
	if err != nil {
		return
	}

	switch jwe.Protected.Enc {
	case JwkEncAlgAes256Gcm:
		c, err := aes.NewCipher(unwrappedKey)
		if err != nil {
			return plaintext, unwrappedKey, err
		}
		gcm, err := cipher.NewGCM(c)
		if err != nil {
			return plaintext, unwrappedKey, err
		}
		ciphertext := make([]byte, len(jwe.Ciphertext)+len(jwe.AuthenticationTag))
		copy(ciphertext, jwe.Ciphertext)
		copy(ciphertext[len(jwe.Ciphertext):], jwe.AuthenticationTag)
		plaintext, err = gcm.Open(plaintext, jwe.InitializationVector, ciphertext, []byte(jwe.Protected.Raw))
		return plaintext, unwrappedKey, err
	default:
		return plaintext, unwrappedKey, fmt.Errorf("unsupported algorithm: %s", jwe.Protected.Enc)
	}

}

func (jwe *JsonWebEncryption) unwrapKey(keyFunc func(header *JoseHeader) crypto.Decrypter) ([]byte, error) {

	decrypter := keyFunc(&jwe.Protected)
	switch jwe.Protected.Alg {
	case JwkEncAlgRsaOeap256:
		return decrypter.Decrypt(nil, jwe.EncryptedKey, &rsa.OAEPOptions{
			Hash: crypto.SHA256,
		})
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", jwe.Protected.Alg)
	}

}
