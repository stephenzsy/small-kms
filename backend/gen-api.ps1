oapi-codegen.exe --package models -generate "types,skip-prune"  ./oapi-api.yaml > models/api_types.gen.go
oapi-codegen.exe --package models -generate echo-server ./oapi-api.yaml > models/api_server.gen.go
oapi-codegen.exe --package agentclient -generate "types,client,skip-prune" -include-tags="agent"  ./oapi-api.yaml  > agent-client/agent_client.gen.go

oapi-codegen.exe --package agentserver -generate "types,echo-server" -include-tags agent ./oapi-agent.yaml  > agent/agentserver/agent_server.gen.go

oapi-codegen.exe --package admin -generate types ./swagger.yaml  > admin/admin_types.gen.go
oapi-codegen.exe --package admin -generate gin-server  ./swagger.yaml  > admin/admin_server.gen.go

oapi-codegen.exe --package client -generate "types,client" -include-tags enroll ./swagger.yaml  > endpoint-enroll/client/enroll_client.gen.go