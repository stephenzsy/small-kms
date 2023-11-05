package cloudkeyaz

import (
	"net/url"
	"strings"
)

func ExtractKeyVaultName(keyvaultEndpoing string) string {
	if parsed, err := url.Parse(keyvaultEndpoing); err == nil {
		return strings.Split(parsed.Host, ".")[0]
	}
	return ""
}
