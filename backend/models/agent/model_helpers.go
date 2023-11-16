package agentmodels

import "github.com/stephenzsy/small-kms/backend/models"

type (
	agentConfigComposed struct {
		models.Ref
		AgentConfigFields
	}
)
