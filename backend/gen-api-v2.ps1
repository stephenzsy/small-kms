oapi-codegen --package models -generate "types,skip-prune" ../api/models-shared.yaml > models/shared_models.gen.go
oapi-codegen --package agentmodels -generate "types,skip-prune" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models" ../api/models-agent.yaml > models/agent/agent_models.gen.go

oapi-codegen --package admin -generate "types,echo-server" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models,models-agent.yaml:github.com/stephenzsy/small-kms/backend/models/agent" ../api/api-v2.yaml > admin/admin.gen.go
