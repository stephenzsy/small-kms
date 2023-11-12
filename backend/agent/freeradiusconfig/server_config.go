package frconfig

import (
	"io"
	"strconv"
	"strings"
)

func (c RadiusServerConfig) MarshalFreeradiusConfig(sb *strings.Builder, linePrefix string) error {
	sb.WriteString("server ")
	sb.WriteString(c.Name)
	sb.WriteString(" {\n")
	FreeRadiusConfigList[RadiusServerListenConfig](c.Listeners).MarshalFreeradiusConfig(sb, "")
	marshalServerAuthorizeConfig(sb, "")
	marshalServerAuthenticateConfig(sb, "")
	sb.WriteString("}\n")
	return nil
}

func (c *RadiusServerConfig) WriteToDigest(digest io.Writer) {
	digest.Write([]byte(c.Name))
	for _, s := range c.Listeners {
		digest.Write([]byte(s.Type))
		digest.Write([]byte(s.Ipaddr))
		digest.Write([]byte(strconv.Itoa(s.Port)))
	}
}

func (c RadiusServerListenConfig) MarshalFreeradiusConfig(sb *strings.Builder, linePrefix string) error {
	sb.WriteString("listen {\n")
	writeStringOmitEmpty(sb, "\t", "type", string(c.Type))
	writeStringOmitEmpty(sb, "\t", "ipaddr", string(c.Ipaddr))
	writeStringOmitEmpty(sb, "\t", "port", strconv.Itoa(c.Port))
	sb.WriteString("}\n")
	return nil
}

var _ FreeRadiusConfigMarshaler = RadiusClientConfig{}

func marshalServerAuthorizeConfig(sb *strings.Builder, linePrefix string) {
	sb.WriteString(linePrefix)
	sb.WriteString("authorize {\n")
	{
		linePrefix := linePrefix + "\t"

		sb.WriteString(linePrefix)
		sb.WriteString("preprocess\n")

		// eap section
		sb.WriteString(linePrefix)
		sb.WriteString("eap {\n")
		{
			linePrefix := linePrefix + "\t"

			sb.WriteString(linePrefix)
			sb.WriteString("ok = return\n")
		}
		sb.WriteString(linePrefix)
		sb.WriteString("}\n")
	}
	sb.WriteString(linePrefix)
	sb.WriteString("}\n")
}

func marshalServerAuthenticateConfig(sb *strings.Builder, linePrefix string) {
	sb.WriteString(linePrefix)
	sb.WriteString("authenticate {\n")

	sb.WriteString(linePrefix)
	sb.WriteString("\teap\n")

	sb.WriteString(linePrefix)
	sb.WriteString("}\n")
}
