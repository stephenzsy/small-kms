package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvKeyAzClientID                 = "AZURE_CLIENT_ID"
	EnvKeyAzTenantID                 = "AZURE_TENANT_ID"
	EnvKeyAzClientSecret             = "AZURE_CLIENT_SECRET"
	EnvKeyAzClientCertPath           = "AZURE_CLIENT_CERTIFICATE_PATH"
	EnvKeyUseManagedIdentity         = "USE_MANAGED_IDENTITY"
	EnvKeyAzKeyvaultResourceEndpoint = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	EnvKeyAzSubscriptionID           = "AZURE_SUBSCRIPTION_ID"
	EnvKeyAzResourceGroupName        = "AZURE_RESOURCE_GROUP_NAME"

	envKeyUseManagedIdentity = "USE_MANAGED_IDENTITY"
)

type EnvService interface {
	Export() []string
	Require(key string, prefixes ...string) (string, bool)
	RequireNonWhitespace(key string, prefixes ...string) (string, bool)
	RequireAbsPath(key string, prefixes ...string) (string, bool)
	Default(key string, value string, prefixes ...string) string

	// convience method to create a error with key missing
	ErrMissing(key string) error
}

type envServiceImpl struct {
	values map[string]*envServiceEntry
}

// RequireAbsPath implements EnvService.
func (impl *envServiceImpl) RequireAbsPath(key string, prefixes ...string) (string, bool) {
	_, entry := impl.resolvePrefixed(key, prefixes)
	p, err := filepath.Abs(entry.raw)
	if err != nil {
		entry.flags &= ^envEntryFlagHasSetValue
		return entry.String(), false
	}
	entry.raw = p
	return entry.String(), entry.HasRequiredValue()
}

// Default implements EnvService.
func (impl *envServiceImpl) Default(key string, value string, prefixes ...string) string {
	_, entry := impl.resolvePrefixed(key, prefixes)
	if entry == nil {
		entry = &envServiceEntry{
			raw: value,
		}
	} else if entry.flags&envEntryFlagHasSetValue == 0 {
		entry.raw = value
	}
	return entry.String()
}

// ErrMissing implements EnvService.
func (*envServiceImpl) ErrMissing(key string) error {
	return fmt.Errorf("missing enviornment variable: %s", key)
}

func (impl *envServiceImpl) Require(key string, prefixes ...string) (string, bool) {
	_, entry := impl.resolvePrefixed(key, prefixes)
	return entry.String(), entry.HasRequiredValue()
}

func (impl *envServiceImpl) RequireNonWhitespace(key string, prefixes ...string) (string, bool) {
	_, entry := impl.resolvePrefixed(key, prefixes)
	if entry != nil {
		entry.raw = strings.TrimSpace(entry.raw)
		if entry.raw == "" {
			entry.flags &= ^envEntryFlagHasSetValue
		}
	}
	return entry.String(), entry.HasRequiredValue()
}

func (impl *envServiceImpl) resolvePrefixed(key string, prefixes []string) (string, *envServiceEntry) {
	for i := range prefixes {
		actualKey := strings.Join(prefixes[i:], "") + key
		if v, ok := impl.lookupEntry(actualKey, false); ok {
			return actualKey, v
		}
	}
	v, _ := impl.lookupEntry(key, true)
	return key, v
}

func (impl *envServiceImpl) lookupEntry(key string, isLeaf bool) (entry *envServiceEntry, ok bool) {
	if v, ok := impl.values[key]; ok {
		return v, ok
	}
	val, ok := os.LookupEnv(key)
	if ok {
		impl.values[key] = &envServiceEntry{
			raw:   val,
			flags: envEntryFlagHasSetValue,
		}
		return impl.values[key], true
	} else if isLeaf {
		impl.values[key] = nil
		return nil, true
	}
	return nil, false
}

func (*envServiceImpl) escapeSlash(s string) string {
	return strings.ReplaceAll(s, "\\", "\\\\")
}

// Export implements EnvService.
func (s *envServiceImpl) Export() []string {
	result := make([]string, 0, len(s.values))
	for key, entry := range s.values {
		result = append(result, fmt.Sprintf("%s=%s", key, s.escapeSlash(entry.String())))
	}
	return result
}

type envServiceEntryFlag uint

const (
	envEntryFlagHasSetValue envServiceEntryFlag = 1 << iota
)

type envServiceEntry struct {
	raw   string
	flags envServiceEntryFlag
}

// HasRequiredValue implements EnvServiceEntry.
func (entry *envServiceEntry) HasRequiredValue() bool {
	if entry == nil {
		return false
	}
	return entry.flags&envEntryFlagHasSetValue != 0
}

// String implements EnvServiceEntry.
func (entry *envServiceEntry) String() string {
	if entry == nil {
		return ""
	}
	return entry.raw
}

var _ EnvService = (*envServiceImpl)(nil)

func NewEnvService() EnvService {
	return &envServiceImpl{
		values: make(map[string]*envServiceEntry),
	}
}
