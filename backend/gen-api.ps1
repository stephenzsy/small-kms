oapi-codegen.exe --package shared -generate "types,skip-prune"  ./oapi-shared.yaml > shared/shared_types.gen.go

oapi-codegen.exe --package base -generate "types,skip-prune" ./oapi-base.yaml > base/base_types.gen.go
oapi-codegen.exe --package base -generate "echo-server" ./oapi-base.yaml > base/base_server.gen.go
oapi-codegen.exe --package profile -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" ./oapi-profile.yaml > profile/profile.gen.go
oapi-codegen.exe --package managedapp -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-profile.yaml:github.com/stephenzsy/small-kms/backend/profile" ./oapi-managed-app.yaml > managedapp/managed_app.gen.go
oapi-codegen.exe --package key -generate "types,echo-server,skip-prune" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" ./oapi-key.yaml > key/key.gen.go
oapi-codegen.exe --package cert -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-key.yaml:github.com/stephenzsy/small-kms/backend/key" ./oapi-cert.yaml > cert/cert.gen.go

oapi-codegen.exe --package agentclient -generate "types,client" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-key.yaml:github.com/stephenzsy/small-kms/backend/key,oapi-cert.yaml:github.com/stephenzsy/small-kms/backend/cert,oapi-managed-app.yaml:github.com/stephenzsy/small-kms/backend/managedapp" -include-tags agentclient ./oapi-client-agent.yaml > agent/client/agent_client.gen.go
oapi-codegen.exe --package agentpush -generate "types" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_types.gen.go
oapi-codegen.exe --package agentpush -generate "client" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_client.gen.go
oapi-codegen.exe --package agentpush -generate "echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_server.gen.go