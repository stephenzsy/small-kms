oapi-codegen --package models -generate "types,skip-prune" ../api/models-shared.yaml > models/shared_models.gen.go
oapi-codegen --package agentmodels -generate "types,skip-prune" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models" ../api/models-agent.yaml > models/agent/agent_models.gen.go
oapi-codegen --package keymodels -generate "types,skip-prune" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models" ../api/models-key.yaml > models/key/key_models.gen.go
oapi-codegen --package certmodels -generate "types,skip-prune" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models,models-key.yaml:github.com/stephenzsy/small-kms/backend/models/key" ../api/models-cert.yaml > models/cert/cert_models.gen.go

oapi-codegen --config "./gen-api-v2-config.yaml" -package admin -generate "models,echo-server" -include-tags "admin" -o "./admin/admin.gen.go" ../api/api-v2.yaml 
oapi-codegen --config "./gen-api-v2-config.yaml" -package agentclient -generate "models,client" -include-tags "agentclient"  -o "./agent/client/v2/agentclient.gen.go" ../api/api-v2.yaml 
oapi-codegen --config "./gen-api-v2-config.yaml" -package agentendpoint -generate "models,client,echo-server" -include-tags "agent"  -o "./agent/endpoint/agent_endpoint.gen.go" ../api/api-v2.yaml 
