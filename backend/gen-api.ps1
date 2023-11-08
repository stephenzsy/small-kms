oapi-codegen --package base -generate "types,skip-prune" ./oapi-base.yaml > base/base_types.gen.go
oapi-codegen --package base -generate "echo-server" ./oapi-base.yaml > base/base_server.gen.go
oapi-codegen --package profile -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" ./oapi-profile.yaml > profile/profile.gen.go
oapi-codegen --package managedapp -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-profile.yaml:github.com/stephenzsy/small-kms/backend/profile,oapi-freeradius-config.yaml:github.com/stephenzsy/small-kms/backend/freeradius/config" ./oapi-managed-app.yaml > managedapp/managed_app.gen.go
oapi-codegen --package key -generate "types,echo-server,skip-prune" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" ./oapi-key.yaml > key/key.gen.go
oapi-codegen --package secret -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base" ./oapi-secret.yaml > secret/secret.gen.go
oapi-codegen --package cert -generate "types,echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-key.yaml:github.com/stephenzsy/small-kms/backend/key" ./oapi-cert.yaml > cert/cert.gen.go

oapi-codegen --package agentclient -generate "types,client" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-key.yaml:github.com/stephenzsy/small-kms/backend/key,oapi-cert.yaml:github.com/stephenzsy/small-kms/backend/cert,oapi-managed-app.yaml:github.com/stephenzsy/small-kms/backend/managedapp" -include-tags agentclient ./oapi-client-agent.yaml > agent/client/agent_client.gen.go
oapi-codegen --package agentpush -generate "types" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-managed-app.yaml:github.com/stephenzsy/small-kms/backend/managedapp" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_types.gen.go
oapi-codegen --package agentpush -generate "client" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-managed-app.yaml:github.com/stephenzsy/small-kms/backend/managedapp" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_client.gen.go
oapi-codegen --package agentpush -generate "echo-server" -import-mapping="oapi-base.yaml:github.com/stephenzsy/small-kms/backend/base,oapi-managed-app.yaml:github.com/stephenzsy/small-kms/backend/managedapp" -include-tags agent ./oapi-agent-push.yaml > agent/push/agent_push_server.gen.go

oapi-codegen --package frconfig -generate "types,skip-prune" ./oapi-freeradius-config.yaml > freeradius/config/freeradius_config_types.gen.go
