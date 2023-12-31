package cert

import (
	"net"
	"net/mail"
	"slices"
	"strings"
)

func SanitizeDNSNames(dnsNames []string) []string {
	for i, v := range dnsNames {
		v = strings.TrimSpace(v)
		v = strings.ToLower(v)
		dnsNames[i] = v
	}
	dnsNames = slices.DeleteFunc(dnsNames, func(s string) bool {
		return s == ""
	})
	slices.Sort(dnsNames)
	dnsNames = slices.Compact(dnsNames)
	if len(dnsNames) == 0 {
		return nil
	}
	return dnsNames
}

func SanitizeEmailAddresses(emailAddresses []string) []string {
	sanitized := make([]string, 0, len(emailAddresses))
	for _, v := range emailAddresses {
		if parsed, err := mail.ParseAddress(v); err == nil {
			sanitized = append(sanitized, parsed.Address)
		}
	}
	if len(sanitized) == 0 {
		return nil
	}
	return sanitized
}

func SanitizeIpAddresses(ips []net.IP) []net.IP {
	ips = slices.DeleteFunc(ips, func(ip net.IP) bool {
		return len(ip) == 0 || ip.IsUnspecified()
	})
	slices.SortFunc(ips, func(a, b net.IP) int {
		lcmp := len(a) - len(b)
		if lcmp != 0 {
			return lcmp
		}
		return slices.Compare(a, b)
	})
	ips = slices.CompactFunc(ips, slices.Equal[net.IP])
	if len(ips) == 0 {
		return nil
	}
	return ips
}

func (sans *SubjectAlternativeNames) Sanitize() *SubjectAlternativeNames {
	if sans == nil {
		return nil
	}
	sans.DNSNames = SanitizeDNSNames(sans.DNSNames)
	sans.Emails = SanitizeEmailAddresses(sans.Emails)
	sans.IPAddresses = SanitizeIpAddresses(sans.IPAddresses)

	if sans.DNSNames == nil && sans.Emails == nil && sans.IPAddresses == nil {
		return nil
	}
	return sans
}
