package secret

import (
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"syscall"
	"unsafe"

	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/secret/internal/ncrypt"
)

var (
	ErrWindowsSyscallFailure = errors.New("winsyscall failed")
	ErrNteKeyExists          = errors.New("key exists")
)

const (
	sizeOfRSABlobHeader = uint32(unsafe.Sizeof(ncrypt.RSAKEY_BLOB{}))
)

type WindowsSecretsService struct {
}

func syscallErrnoIs(err error, expected syscall.Errno) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == expected
	}
	return false
}

// SignData implements SecretService.
func (s *WindowsSecretsService) RS256SignHash(hash []byte, certName string) ([]byte, *rsa.PublicKey, error) {
	badSyscall := func(e error, locator string) ([]byte, *rsa.PublicKey, error) {
		return nil, nil, fmt.Errorf("%w:%s:%w", ErrWindowsSyscallFailure, locator, e)
	}

	// open provider
	hProvider, err := s.openStorageProvider()
	if err != nil {
		return badSyscall(err, "openprovider")
	}
	defer ncrypt.FreeObject(ncrypt.HANDLE(hProvider))

	// open key
	hKey, err := s.getPersistedRSAKey(hProvider, 2048, certName)
	if err != nil {
		return badSyscall(err, "openkey")
	}
	defer ncrypt.FreeObject(ncrypt.HANDLE(hKey))

	// sign
	info := ncrypt.PKCS1_PADDING_INFO{
		AlgId: utf16PtrFromString(ncrypt.SHA256_ALGORITHM),
	}
	var signatureSize uint32
	err = ncrypt.SignHash(hKey, unsafe.Pointer(&info), hash, nil, &signatureSize, ncrypt.PAD_PKCS1)
	if err != nil {
		return badSyscall(err, "sign-len")
	}
	signature := make([]byte, signatureSize)
	err = ncrypt.SignHash(hKey, unsafe.Pointer(&info), hash, signature, &signatureSize, ncrypt.PAD_PKCS1)
	if err != nil {
		return badSyscall(err, "sign")
	}

	// export public key
	var pubkeySize uint32
	psBlobType := utf16PtrFromString(ncrypt.RSAPUBLIC_KEY_BLOB)
	if err = ncrypt.ExportKey(hKey, 0, psBlobType, nil, nil, &pubkeySize, 0); err != nil {
		return badSyscall(err, "exportkey-len")
	}
	keyblobBytes := make([]byte, pubkeySize)
	if err = ncrypt.ExportKey(hKey, 0, psBlobType, nil, keyblobBytes, &pubkeySize, 0); err != nil {
		return badSyscall(err, "exportkey")
	}
	keyBlob, data, err := convertWinPublicKeyBlob(keyblobBytes)
	if err != nil {
		return nil, nil, err
	}
	consumeBigInt := func(size uint32) *big.Int {
		b := data[:size]
		data = data[size:]
		return new(big.Int).SetBytes(b)
	}
	pk := rsa.PublicKey{
		E: int(consumeBigInt(keyBlob.PublicExpSize).Int64()),
		N: consumeBigInt(keyBlob.ModulusSize),
	}
	return signature[:signatureSize], &pk, err
}

func convertWinPublicKeyBlob(blob []byte) (ncrypt.RSAKEY_BLOB, []byte, error) {
	if len(blob) < int(sizeOfRSABlobHeader) {
		return ncrypt.RSAKEY_BLOB{}, nil, errors.New("cng: exported key is corrupted")
	}
	hdr := (*(*ncrypt.RSAKEY_BLOB)(unsafe.Pointer(&blob[0])))
	return hdr, blob[sizeOfRSABlobHeader:], nil
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

// GenerateRSAKey implements SecretsService.
func (s *WindowsSecretsService) getPersistedRSAKey(hProvider ncrypt.PROV_HANDLE, keyLength int, keyName string) (hKey ncrypt.KEY_HANDLE, err error) {

	var keyNamePtr *uint16
	if keyName != "" {
		keyNamePtr = utf16PtrFromString(keyName)
	}
	err = ncrypt.CreatePersistedKey(hProvider, &hKey,
		utf16PtrFromString(ncrypt.RSA_ALGORITHM),
		keyNamePtr,
		0,
		0) //ncrypt.NCRYPT_MACHINE_KEY_FLAG)
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

func (s *WindowsSecretsService) openStorageProvider() (hProvider ncrypt.PROV_HANDLE, err error) {
	err = ncrypt.OpenStorageProvider(&hProvider, nil, 0)
	if err != nil {
		err = fmt.Errorf("%w:%w", ErrWindowsSyscallFailure, err)
	}
	return
}
