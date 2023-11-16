package agentmodels

import "github.com/stephenzsy/small-kms/backend/models"

type (
	agentComposed struct {
		models.ApplicationByAppId
		AgentFields
	}
)
