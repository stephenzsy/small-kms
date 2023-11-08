#!/usr/bin/env pwsh
$env:GOARCH = "arm64"
$env:GOOS = "darwin"
go build -o smallkms.out . && ./smallkms.out --env-file ./.env --debug --pretty-log admin "localhost:9001"
