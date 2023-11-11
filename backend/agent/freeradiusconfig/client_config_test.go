package frconfig_test

import (
	"strings"
	"testing"

	frconfig "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
)

func TestClientConfig_MarshalText(t *testing.T) {
	c := &frconfig.RadiusClientConfig{
		Name:   "localhost",
		Ipaddr: "127.0.0.1",
		Secret: "testing123",
	}

	sb := &strings.Builder{}
	err := c.MarshalFreeradiusConfig(sb, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := "client localhost {\n\tipaddr = 127.0.0.1\n\tsecret = testing123\n}\n"
	if sb.String() != expected {
		t.Errorf("expected %q, but got %q", expected, sb.String())
	}
}
