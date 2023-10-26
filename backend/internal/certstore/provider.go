package certstore

import "crypto"

type KeySession interface {
	crypto.Signer
	MarkKeyPersistent()
	Close()
}

type CryptoStoreProvider interface {
	CreateRSAKeySession(keyName string, keyLength int, isMachineLevel bool) (ks KeySession, err error)
}
