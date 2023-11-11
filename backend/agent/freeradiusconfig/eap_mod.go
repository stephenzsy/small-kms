package frconfig

import "strings"

func MarshalModEAP(sb *strings.Builder, linePrefix string) error {
	sb.WriteString("eap {\n")
	writeStringOmitEmpty(sb, "\t", "default_eap_type", "tls")
	writeStringOmitEmpty(sb, "\t", "timer_expire", "60")
	writeStringOmitEmpty(sb, "\t", "ignore_unknown_eap_types", "no")
	writeStringOmitEmpty(sb, "\t", "cisco_accounting_username_bug", "no")
	writeStringOmitEmpty(sb, "\t", "max_sessions", "${max_requests}")
	MarshalRadiusEapTlsConfig(sb, "\t")
	sb.WriteString("}\n")
	return nil
}

func MarshalRadiusEapTlsConfig(sb *strings.Builder, linePrefix string) error {
	innerLinePrefix := linePrefix + "\t"

	sb.WriteString(linePrefix)
	sb.WriteString("tls-config tls-common {\n")

	writeStringOmitEmpty(sb, innerLinePrefix, "private_key_file", "${certdir}/server.pem")
	writeStringOmitEmpty(sb, innerLinePrefix, "certificate_file", "${certdir}/server.pem")
	writeStringOmitEmpty(sb, innerLinePrefix, "ca_file", "${cadir}/ca.pem")
	writeStringOmitEmpty(sb, innerLinePrefix, "cipher_list", "\"DEFAULT\"")
	writeStringOmitEmpty(sb, innerLinePrefix, "cipher_server_preference", "yes")
	writeStringOmitEmpty(sb, innerLinePrefix, "tls_min_version", "1.2")
	writeStringOmitEmpty(sb, innerLinePrefix, "tls_max_version", "1.3")
	writeStringOmitEmpty(sb, innerLinePrefix, "ecdh_curve", "secp384r1")

	sb.WriteString(linePrefix)
	sb.WriteString("}\n")

	sb.WriteString(linePrefix)
	sb.WriteString("tls {\n")

	writeStringOmitEmpty(sb, innerLinePrefix, "tls", "tls-common")

	sb.WriteString(linePrefix)
	sb.WriteString("}\n")

	return nil
}
