oapi-codegen.exe --package shared -generate "types,skip-prune"  ./oapi-shared.yaml > shared/shared_types.gen.go

oapi-codegen.exe --package models -generate "types" -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared"  ./oapi-api.yaml > models/api_types.gen.go
oapi-codegen.exe --package models -generate echo-server -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" -include-tags="admin,diagnostics" ./oapi-api.yaml > models/api_server.gen.go
oapi-codegen.exe --package agentclient -generate "types,client" -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" -include-tags="agentclient,diagnostics" ./oapi-api.yaml  > agent-client/agent_client.gen.go

oapi-codegen.exe --package agentserver -generate "types,echo-server" -include-tags agent -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" ./oapi-api.yaml > agent/agentserver/agent_server.gen.go
oapi-codegen.exe --package agentproxyclient -generate "types,client" -include-tags agentproxyclient -import-mapping="oapi-shared.yaml:github.com/stephenzsy/small-kms/backend/shared" ./oapi-api.yaml > admin/agentproxyclient/agent_proxy_client.gen.go

#oapi-codegen.exe --package client -generate "types,client" -include-tags enroll ./swagger.yaml  > endpoint-enroll/client/enroll_client.gen.go

oapi-codegen.exe --package base -generate "types,skip-prune" ./oapi-base.yaml > base/base_types.gen.go
oapi-codegen.exe --package managedapp -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" -include-tags="admin" ./oapi-managed-app.yaml > managedapp/managed_app.gen.go
