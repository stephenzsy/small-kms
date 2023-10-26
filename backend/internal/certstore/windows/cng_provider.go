package wincryptostore

import (
	"crypto"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"syscall"
	"unsafe"

	"github.com/microsoft/go-crypto-winnative/cng"
	"github.com/stephenzsy/small-kms/backend/internal/certstore"
	"github.com/stephenzsy/small-kms/backend/internal/certstore/windows/ncrypt"
)

var (
	ErrWindowsSyscallFailure = errors.New("winsyscall failed")
	ErrNteKeyExists          = errors.New("key exists")
)

func syscallErrnoIs(err error, expected syscall.Errno) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == expected
	}
	return false
}

type windowsNCryptKeySession struct {
	publicKey crypto.PublicKey
	hProvider ncrypt.PROV_HANDLE
	hKey      ncrypt.KEY_HANDLE
	toPersist bool
}

// Public implements crypto.Signer.
func (ks *windowsNCryptKeySession) Public() crypto.PublicKey {
	return ks.publicKey
}

// Sign implements crypto.Signer.
func (ks *windowsNCryptKeySession) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	badSyscall := func(e error, locator string) ([]byte, error) {
		return nil, fmt.Errorf("%w:%s:%w", ErrWindowsSyscallFailure, locator, e)
	}

	var algId *uint16
	switch opts.HashFunc() {
	case crypto.SHA256:
		algId = utf16PtrFromString(ncrypt.SHA256_ALGORITHM)
	default:
		return nil, fmt.Errorf("unsupported hash function: %v", opts.HashFunc())
	}
	info := ncrypt.PKCS1_PADDING_INFO{
		AlgId: algId,
	}
	var signatureSize uint32
	err = ncrypt.SignHash(ks.hKey, unsafe.Pointer(&info), digest, nil, &signatureSize, ncrypt.PAD_PKCS1)
	if err != nil {
		return badSyscall(err, "sign-len")
	}
	signature = make([]byte, signatureSize)
	err = ncrypt.SignHash(ks.hKey, unsafe.Pointer(&info), digest, signature, &signatureSize, ncrypt.PAD_PKCS1)
	if err != nil {
		return badSyscall(err, "sign")
	}
	return signature, nil
}

// Close implements certstore.KeySession.
func (ks *windowsNCryptKeySession) Close() {
	if !ks.toPersist {
		ncrypt.DeleteKey(ks.hKey, ncrypt.NCRYPT_SILENT_FLAG)
	}
	if ks.hKey != 0 {
		ncrypt.FreeObject(ncrypt.HANDLE(ks.hKey))
	}
	if ks.hProvider != 0 {
		ncrypt.FreeObject(ncrypt.HANDLE(ks.hProvider))
	}
}

// MarkKeyPersistent implements certstore.KeySession.
func (ks *windowsNCryptKeySession) MarkKeyPersistent() {
	ks.toPersist = true
}

var _ certstore.KeySession = (*windowsNCryptKeySession)(nil)

type windowsNCryptCryptoStoreProvider struct {
}

// GenerateRSAKeyPair implements certstore.CryptoStoreProvider.
func (*windowsNCryptCryptoStoreProvider) GenerateRSAKeyPair(keyLength int) (*rsa.PrivateKey, error) {
	n, e, d, p, q, dp, dq, qi, err := cng.GenerateKeyRSA(keyLength)
	if err != nil {
		return nil, err
	}
	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Uint64()),
		},
		D: new(big.Int).SetBytes(d),
		Primes: []*big.Int{
			new(big.Int).SetBytes(p),
			new(big.Int).SetBytes(q),
		},
		Precomputed: rsa.PrecomputedValues{
			Dp:   new(big.Int).SetBytes(dp),
			Dq:   new(big.Int).SetBytes(dq),
			Qinv: new(big.Int).SetBytes(qi),
		},
	}, nil
}

func utf16PtrFromString(s string) *uint16 {
	return &utf16FromString(s)[0]
}

// utf16FromString converts the string using a stack-allocated slice of 64 bytes.
// It should only be used to convert known BCrypt identifiers which only contains ASCII characters.
// utf16FromString allocates if s is longer than 31 characters.
func utf16FromString(s string) []uint16 {
	// Once https://go.dev/issues/51896 lands and our support matrix allows it,
	// we can replace part of this function by utf16.AppendRune
	a := make([]uint16, 0, 32)
	for _, v := range s {
		if v == 0 || v > 127 {
			panic("utf16FromString only supports ASCII characters, got " + s)
		}
		a = append(a, uint16(v))
	}
	// Finish with a NULL byte.
	a = append(a, 0)
	return a
}

const (
	sizeOfRSABlobHeader = uint32(unsafe.Sizeof(ncrypt.RSAKEY_BLOB{}))
)

func convertWinPublicKeyBlob(blob []byte) (ncrypt.RSAKEY_BLOB, []byte, error) {
	if len(blob) < int(sizeOfRSABlobHeader) {
		return ncrypt.RSAKEY_BLOB{}, nil, errors.New("cng: exported key is corrupted")
	}
	hdr := (*(*ncrypt.RSAKEY_BLOB)(unsafe.Pointer(&blob[0])))
	return hdr, blob[sizeOfRSABlobHeader:], nil
}

func (p *windowsNCryptCryptoStoreProvider) CreateRSAKeySession(keyName string, keyLength int, isMachineLevel bool) (certstore.KeySession, error) {
	ks := &windowsNCryptKeySession{}
	var err error
	badSyscall := func(e error, locator string) (*windowsNCryptKeySession, error) {
		return ks, fmt.Errorf("%w:%s:%w", ErrWindowsSyscallFailure, locator, e)
	}

	// open provider
	ks.hProvider, err = p.openStorageProvider()
	if err != nil {
		return badSyscall(err, "openprovider")
	}

	// open key
	//ncrypt.NCRYPT_MACHINE_KEY_FLAG)
	flag := uint32(0)
	if isMachineLevel {
		flag = ncrypt.NCRYPT_MACHINE_KEY_FLAG
	}

	ks.hKey, err = p.getPersistedRSAKey(ks.hProvider, keyLength, keyName, flag)
	if err != nil {
		return badSyscall(err, "openkey")
	}

	// export public key
	var pubkeySize uint32
	psBlobType := utf16PtrFromString(ncrypt.RSAPUBLIC_KEY_BLOB)
	if err = ncrypt.ExportKey(ks.hKey, 0, psBlobType, nil, nil, &pubkeySize, 0); err != nil {
		return badSyscall(err, "exportkey-len")
	}
	keyblobBytes := make([]byte, pubkeySize)
	if err = ncrypt.ExportKey(ks.hKey, 0, psBlobType, nil, keyblobBytes, &pubkeySize, 0); err != nil {
		return badSyscall(err, "exportkey")
	}
	keyBlob, data, err := convertWinPublicKeyBlob(keyblobBytes)
	if err != nil {
		return ks, err
	}
	consumeBigInt := func(size uint32) *big.Int {
		b := data[:size]
		data = data[size:]
		return new(big.Int).SetBytes(b)
	}
	ks.publicKey = &rsa.PublicKey{
		E: int(consumeBigInt(keyBlob.PublicExpSize).Int64()),
		N: consumeBigInt(keyBlob.ModulusSize),
	}
	return ks, nil
}

func (s *windowsNCryptCryptoStoreProvider) openStorageProvider() (hProvider ncrypt.PROV_HANDLE, err error) {
	err = ncrypt.OpenStorageProvider(&hProvider, nil, 0)
	if err != nil {
		err = fmt.Errorf("%w:%w", ErrWindowsSyscallFailure, err)
	}
	return
}

// GenerateRSAKey implements SecretsService.
func (s *windowsNCryptCryptoStoreProvider) getPersistedRSAKey(
	hProvider ncrypt.PROV_HANDLE, keyLength int, keyName string, createKeyFlag uint32) (hKey ncrypt.KEY_HANDLE, err error) {

	var keyNamePtr *uint16
	if keyName != "" {
		keyNamePtr = utf16PtrFromString(keyName)
	}
	err = ncrypt.CreatePersistedKey(hProvider, &hKey,
		utf16PtrFromString(ncrypt.RSA_ALGORITHM),
		keyNamePtr,
		0,
		createKeyFlag)
	if err != nil {
		if syscallErrnoIs(err, ncrypt.NTE_EXISTS) {
			err = ncrypt.OpenKey(hProvider, &hKey, keyNamePtr, 0, 0)
			return
		}
		if err != nil {
			return
		}
	}

	// set key length
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(keyLength))
	err = ncrypt.SetProperty(
		ncrypt.HANDLE(hKey),
		utf16PtrFromString(ncrypt.NCRYPT_LENGTH_PROPERTY),
		&bs[0],
		uint32(len(bs)),
		0)
	if err != nil {
		return
	}

	exportPolicy := make([]byte, 4)
	binary.LittleEndian.PutUint32(exportPolicy, uint32(0))
	err = ncrypt.SetProperty(
		ncrypt.HANDLE(hKey),
		utf16PtrFromString(ncrypt.NCRYPT_EXPORT_POLICY_PROPERTY),
		&exportPolicy[0],
		uint32(len(exportPolicy)),
		0)
	if err != nil {
		return
	}

	err = ncrypt.FinalizeKey(hKey, 0)
	if err != nil {
		return
	}
	return
}

var _ certstore.CryptoStoreProvider = (*windowsNCryptCryptoStoreProvider)(nil)

func NewWindowsNCryptCryptoStoreProvider() *windowsNCryptCryptoStoreProvider {
	return &windowsNCryptCryptoStoreProvider{}
}
