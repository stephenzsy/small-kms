 oapi-codegen.exe --package admin -generate types ./swagger.yaml  > admin/admin_types.gen.go
 oapi-codegen.exe --package admin -generate gin-server ./swagger.yaml  > admin/admin_server.gen.go
 oapi-codegen.exe --package common -generate types,skip-prune ./enums.openapi.yaml  > common/enums.gen.go

 #oapi-codegen.exe --package admin -generate client ./swagger.yaml  > admin/admin_client.gen.go