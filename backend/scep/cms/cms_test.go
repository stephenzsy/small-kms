package cms

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestParseSamplePkiPayload(t *testing.T) {
	b, err := os.ReadFile("testdata/pki-payload.txt")
	if err != nil {
		t.Skip(err)
	}
	payload, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		t.Error(err)
	}
	_, err = ParsePkiMessage(payload)
	if err != nil {
		t.Error(err)
	}
}
