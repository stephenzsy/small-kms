oapi-codegen.exe --package shared -generate "types,skip-prune"  ./oapi-shared.yaml > shared/shared_types.gen.go

oapi-codegen.exe --package models -generate "types" -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared"  ./oapi-api.yaml > models/api_types.gen.go
oapi-codegen.exe --package models -generate echo-server -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" ./oapi-api.yaml > models/api_server.gen.go
oapi-codegen.exe --package agentclient -generate "types,client" -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" -include-tags="agent" ./oapi-api.yaml  > agent-client/agent_client.gen.go

oapi-codegen.exe --package agentserver -generate "types,echo-server" -include-tags agent ./oapi-agent.yaml  > agent/agentserver/agent_server.gen.go

oapi-codegen.exe --package admin -generate types ./swagger.yaml  > admin/admin_types.gen.go
oapi-codegen.exe --package admin -generate gin-server  ./swagger.yaml  > admin/admin_server.gen.go

oapi-codegen.exe --package client -generate "types,client" -include-tags enroll ./swagger.yaml  > endpoint-enroll/client/enroll_client.gen.go