oapi-codegen --package models -generate "types,skip-prune" ../api/models-shared.yaml > models/models_shared.gen.go
oapi-codegen --package admin -generate "types,echo-server" -import-mapping="models-shared.yaml:github.com/stephenzsy/small-kms/backend/models" ../api/api-v2.yaml > admin/admin.gen.go
