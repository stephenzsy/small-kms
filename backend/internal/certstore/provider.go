package certstore

import (
	"crypto"
	"crypto/rsa"
)

type KeySession interface {
	crypto.Signer
	MarkKeyPersistent()
	Close()
}

type CryptoStoreProvider interface {
	CreateRSAKeySession(keyName string, keyLength int, isMachineLevel bool) (ks KeySession, err error)
	GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error)
}
