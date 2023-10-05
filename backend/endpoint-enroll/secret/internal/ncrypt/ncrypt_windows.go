package ncrypt

import "syscall"

type (
	HANDLE      syscall.Handle
	PROV_HANDLE HANDLE
	KEY_HANDLE  HANDLE
)

// error codes
const (
	NTE_EXISTS = 0x8009000F
)

// https://learn.microsoft.com/en-us/windows/win32/seccng/key-storage-property-identifiers
const (
	NCRYPT_PERSIST_FLAG       = 0x00000000
	NCRYPT_SILENT_FLAG        = 0x00000040
	NCRYPT_MACHINE_KEY_FLAG   = 0x00000020
	NCRYPT_OVERWRITE_KEY_FLAG = 0x00000080
	NCRYPT_PERSIST_ONLY_FLAG  = 0x40000000
)

const (
	RSAPUBLIC_KEY_BLOB = "RSAPUBLICBLOB"
)

const (
	NCRYPT_LENGTH_PROPERTY        = "Length"
	NCRYPT_EXPORT_POLICY_PROPERTY = "Export Policy"
)

// https://learn.microsoft.com/en-us/windows/win32/com/com-error-codes-4

// Selected list of https://learn.microsoft.com/en-us/windows/win32/seccng/cng-algorithm-identifiers used by this application
const (
	RSA_ALGORITHM   = "RSA"
	ECDSA_ALGORITHM = "ECDSA"
)

const (
	SHA256_ALGORITHM = "SHA256"
)

// https://learn.microsoft.com/en-us/windows/win32/api/ncrypt/nf-ncrypt-ncryptopenstorageprovider
const (
	MS_KEY_STORAGE_PROVIDER            = "Microsoft Software Key Storage Provider"
	MS_SMART_CARD_KEY_STORAGE_PROVIDER = "Microsoft Smart Card Key Storage Provider"
	MS_PLATFORM_CRYPTO_PROVIDER        = "Microsoft Platform Crypto Provider"
)

// https://learn.microsoft.com/en-us/windows/win32/seccng/key-storage-property-identifiers
const (
	NCRYPT_ALLOW_EXPORT_FLAG              = 0x00000001
	NCRYPT_ALLOW_PLAINTEXT_EXPORT_FLAG    = 0x00000002
	NCRYPT_ALLOW_ARCHIVING_FLAG           = 0x00000004
	NCRYPT_ALLOW_PLAINTEXT_ARCHIVING_FLAG = 0x00000008
)

type KeyBlobMagicNumber uint32
type PadMode uint32

// https://docs.microsoft.com/en-us/windows/win32/api/bcrypt/ns-bcrypt-bcrypt_rsakey_blob
type RSAKEY_BLOB struct {
	Magic         KeyBlobMagicNumber
	BitLength     uint32
	PublicExpSize uint32
	ModulusSize   uint32
	Prime1Size    uint32
	Prime2Size    uint32
}

// https://docs.microsoft.com/en-us/windows/win32/api/bcrypt/ns-bcrypt-bcrypt_ecckey_blob
type ECCKEY_BLOB struct {
	Magic   KeyBlobMagicNumber
	KeySize uint32
}

const (
	PAD_UNDEFINED PadMode = 0x0
	PAD_NONE      PadMode = 0x1
	PAD_PKCS1     PadMode = 0x2
	PAD_OAEP      PadMode = 0x4
	PAD_PSS       PadMode = 0x8
)

type PKCS1_PADDING_INFO struct {
	AlgId *uint16
}

type Buffer struct {
	Length uint32
	Type   uint32
	Data   uintptr
}

type BufferDesc struct {
	Version uint32
	Count   uint32 // number of buffers
	Buffers *Buffer
}

//sys   CreatePersistedKey(hProvider PROV_HANDLE, phKey *KEY_HANDLE, pszAlgId *uint16, pszKeyName *uint16, dwLegacyKeySpec uint32, dwFlags uint32) (s error) = ncrypt.NCryptCreatePersistedKey
//sys   ExportKey(hKey KEY_HANDLE, hExportKey KEY_HANDLE, pszBlobType *uint16, pParameterList *BufferDesc, pbOutput []byte, pcbResult *uint32, dwFlags uint32) (s error) = ncrypt.NCryptExportKey
//sys   FinalizeKey(hKey KEY_HANDLE, dwFlags uint32) (s error) = ncrypt.NCryptFinalizeKey
//sys   FreeObject(hObject HANDLE) (s error) = ncrypt.NCryptFreeObject
//sys   OpenKey(hProvider PROV_HANDLE, phKey *KEY_HANDLE, pszKeyName *uint16, dwLegacyKeySpec uint32, dwFlags uint32) (s error) = ncrypt.NCryptOpenKey
//sys   OpenStorageProvider(phProvider *PROV_HANDLE, pszProviderName *uint16, dwFlags uint32) (s error) = ncrypt.NCryptOpenStorageProvider
//sys   SetProperty(hObject HANDLE, pszProperty *uint16, pbInput *byte, cbInput uint32, dwFlags uint32) (s error) = ncrypt.NCryptSetProperty
//sys   SignHash (hKey KEY_HANDLE, pPaddingInfo unsafe.Pointer, pbHashValue []byte, pbSignature []byte, pcbResult *uint32, dwFlags PadMode) (s error) = ncrypt.NCryptSignHash
