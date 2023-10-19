package agentconfig

import (
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentCallbackDoc struct {
	kmsdoc.BaseDoc

	Name    shared.AgentConfigName   `json:"name"`    // for index only
	Version kmsdoc.HexStringStroable `json:"version"` // version of the config after evaluation
}
