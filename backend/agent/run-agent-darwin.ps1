#!/usr/bin/env pwsh
$env:GOARCH = "arm64"
$env:GOOS = "darwin"
go build -ldflags="-X 'main.BuildID=dev-test'" -o agent.exe . && ./agent.exe --pretty-log -env .env server :10443
