package frconfig

import (
	"encoding"
	"strings"
)

func (*ClientConfig) GetType() string {
	return "client"
}

func (c *ClientConfig) GetName() string {
	return c.Name
}

func writeStringOmitEmpty(sb *strings.Builder, prefixTabs string, key string, value string) {
	if value != "" {
		sb.WriteString(prefixTabs)
		sb.WriteString(key)
		sb.WriteString(" = ")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
}

func (c ClientConfig) MarshalText() ([]byte, error) {
	sb := &strings.Builder{}
	sb.WriteString(c.GetType())
	sb.WriteString(" ")
	sb.WriteString(c.GetName())
	sb.WriteString(" {\n")
	writeStringOmitEmpty(sb, "\t", "ipaddr", c.Ipaddr)
	writeStringOmitEmpty(sb, "\t", "secret", c.Secret)
	sb.WriteString("}\n")
	return []byte(sb.String()), nil
}

var _ encoding.TextMarshaler = (*ClientConfig)(nil)
