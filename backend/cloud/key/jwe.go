package cloudkey

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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
	JwkEncAlgEcdhEs     JsonWebKeyEncryptionAlgorithm = "ECDH-ES"
)

type JoseHeader struct {
	Algorithm           JsonWebKeyEncryptionAlgorithm `json:"alg,omitempty"`
	EncryptionAlgorithm JsonWebKeyEncryptionAlgorithm `json:"enc"`
	KeyID               string                        `json:"kid,omitempty"`
	EphemeralPublicKey  *JsonWebKey                   `json:"epk,omitempty"`
	AgreementPartyUInfo Base64RawURLEncodableBytes    `json:"apu,omitempty"`
	AgreementPartyVInfo Base64RawURLEncodableBytes    `json:"apv,omitempty"`
	Raw                 string                        `json:"-"`
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

func (jwe *JsonWebEncryption) Decrypt(keyFunc func(header *JoseHeader) crypto.PrivateKey) (plaintext []byte, unwrappedKey []byte, err error) {

	unwrappedKey, err = jwe.unwrapKey(keyFunc)
	if err != nil {
		return
	}

	switch jwe.Protected.EncryptionAlgorithm {
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
		return plaintext, unwrappedKey, fmt.Errorf("unsupported algorithm: %s", jwe.Protected.EncryptionAlgorithm)
	}

}

func (jwe *JsonWebEncryption) unwrapKey(keyFunc func(header *JoseHeader) crypto.PrivateKey) ([]byte, error) {

	privateKey := keyFunc(&jwe.Protected)
	switch jwe.Protected.Algorithm {
	case JwkEncAlgRsaOeap256:
		if privateKey, ok := privateKey.(crypto.Decrypter); !ok {
			return nil, fmt.Errorf("incompatable key")
		} else {
			return privateKey.Decrypt(nil, jwe.EncryptedKey, &rsa.OAEPOptions{
				Hash: crypto.SHA256,
			})
		}
	case JwkEncAlgEcdhEs:
		if jwe.Protected.EncryptionAlgorithm != JwkEncAlgAes256Gcm {
			return nil, fmt.Errorf("incompatable enc")
		}
		if privateKey, ok := privateKey.(*ecdh.PrivateKey); !ok {
			return nil, fmt.Errorf("incompatable key")
		} else if epk, ok := jwe.Protected.EphemeralPublicKey.PublicKey().(*ecdsa.PublicKey); !ok {
			return nil, fmt.Errorf("incompatable key, key should be ecdsa")
		} else if ecdhPubKey, err := epk.ECDH(); err != nil {
			return nil, err
		} else if z, err := privateKey.ECDH(ecdhPubKey); err != nil {
			return nil, err
		} else {
			kdf := &ecdhesKDF{
				z:   z[:32],
				alg: string(jwe.Protected.EncryptionAlgorithm),
				apu: jwe.Protected.AgreementPartyUInfo,
				apv: jwe.Protected.AgreementPartyVInfo,
			}
			return kdf.getAESGCM256DerivedKey(), nil
		}
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", jwe.Protected.Algorithm)
	}

}

type ecdhesKDF struct {
	z   []byte
	alg string
	apu []byte
	apv []byte
}

func uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

func (kdf *ecdhesKDF) getAESGCM256DerivedKey() []byte {
	d := sha256.New()
	d.Write(uint32ToBytes(1)) // round 1
	d.Write(kdf.z)
	d.Write(uint32ToBytes(uint32(len(JwkEncAlgAes256Gcm))))
	d.Write([]byte(JwkEncAlgAes256Gcm))
	d.Write(uint32ToBytes(uint32(len(kdf.apu))))
	d.Write(kdf.apu)
	d.Write(uint32ToBytes(uint32(len(kdf.apv))))
	d.Write(kdf.apv)
	d.Write(uint32ToBytes(256))
	return d.Sum(nil)
}
