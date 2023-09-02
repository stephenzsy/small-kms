 oapi-codegen.exe --package admin -generate types,gin-server,skip-prune ./swagger.yaml  > admin/admin_server.gen.go
 oapi-codegen.exe --package scep -generate types,gin-server,skip-prune ./scep/swagger.yaml  > scep/scep_server.gen.go
 oapi-codegen.exe --package msintune -generate types,client ./scep/msintune/swagger.yaml  > scep/msintune/msintune_client.gen.go