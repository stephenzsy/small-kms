package radius

import (
	"strings"
)

func writeStringOmitEmpty(sb *strings.Builder, prefixTabs string, key string, value string) {
	if value != "" {
		sb.WriteString(prefixTabs)
		sb.WriteString(key)
		sb.WriteString(" = ")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
}

func (c RadiusClientConfig) MarshalFreeradiusConfig() ([]byte, error) {
	sb := &strings.Builder{}
	sb.WriteString("client ")
	sb.WriteString(c.Name)
	sb.WriteString(" {\n")
	writeStringOmitEmpty(sb, "\t", "ipaddr", c.Ipaddr)
	writeStringOmitEmpty(sb, "\t", "secret", c.Secret)
	sb.WriteString("}\n")
	return []byte(sb.String()), nil
}
