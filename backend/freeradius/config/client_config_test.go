package frconfig_test

import (
	"testing"

	frconfig "github.com/stephenzsy/small-kms/backend/freeradius/config"
)

func TestClientConfig_MarshalText(t *testing.T) {
	c := &frconfig.RadiusClientConfig{
		Name:   "localhost",
		Ipaddr: "127.0.0.1",
		Secret: "testing123",
	}

	b, err := c.MarshalFreeradiusConfig()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := "client localhost {\n\tipaddr = 127.0.0.1\n\tsecret = testing123\n}\n"
	if string(b) != expected {
		t.Errorf("expected %q, but got %q", expected, string(b))
	}
}
