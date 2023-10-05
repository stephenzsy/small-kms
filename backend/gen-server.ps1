oapi-codegen.exe --package models -generate "types,skip-prune" ./oapi-api.yml > models/api_types.gen.go
oapi-codegen.exe --package models -generate gin-server ./oapi-api.yml > models/api.gen.go
oapi-codegen.exe --package admin -generate types -import-mapping="oapi-api.yml:github.com/stephenzsy/small-kms/backend/models" ./swagger.yaml  > admin/admin_types.gen.go
oapi-codegen.exe --package admin -generate gin-server -import-mapping="oapi-api.yml:github.com/stephenzsy/small-kms/backend/models" ./swagger.yaml  > admin/admin_server.gen.go
oapi-codegen.exe --package client -generate "types,client" -import-mapping="oapi-api.yml:github.com/stephenzsy/small-kms/backend/models" -include-tags enroll ./swagger.yaml  > endpoint-enroll/client/enroll_client.gen.go