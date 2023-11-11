package radius

import "github.com/stephenzsy/small-kms/backend/managedapp"

type AgentConfigRadius = managedapp.AgentConfigRadius

type contextKey string

const (
	contextKeyRadiusConfig          contextKey = "radiusConfig"
	contextKeyRadiusConfigProcessed contextKey = "radiusConfigProcessed"
)
