$env:GOARCH="amd64"
$env:GOOS="linux"
go build  -o ./agent/agent ./agent

