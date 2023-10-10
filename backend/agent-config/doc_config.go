package agentconfig

import "github.com/stephenzsy/small-kms/backend/internal/kmsdoc"

type AgentConfigDoc struct {
	kmsdoc.BaseDoc

	Name    string                   `json:"name"`    // for index only
	Version kmsdoc.HexStringStroable `json:"version"` // md5 checksum of the config
}
