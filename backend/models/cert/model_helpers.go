package certmodels

import (
	"crypto/x509/pkix"
	"io"
	"net"
	"net/mail"
	"slices"
	"strings"

	"github.com/stephenzsy/small-kms/backend/models"
)

type (
	certificatePolicyComposed struct {
		models.Ref
		CertificatePolicyFields
	}

	certificateRefComposed struct {
		models.Ref
		CertificateRefFields
	}

	certificateComposed struct {
		CertificateRef
		CertificateFields
	}
)

func (cs *CertificateSubject) String() string {
	return cs.ToPkixName().String()
}

func (cs *CertificateSubject) ToPkixName() pkix.Name {
	return pkix.Name{
		CommonName: cs.CommonName,
	}
}

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

func (sans *SubjectAlternativeNames) Digest(w io.Writer) {
	if sans == nil {
		return
	}
	for _, v := range sans.DNSNames {
		io.WriteString(w, v)
	}
	for _, v := range sans.Emails {
		io.WriteString(w, v)
	}
	for _, v := range sans.IPAddresses {
		w.Write(v)
	}
}
