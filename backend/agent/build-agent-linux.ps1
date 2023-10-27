$env:GOARCH = "amd64"
$env:GOOS = "linux"
$env:CGO_ENABLED = 1
go build -ldflags="-X 'main.BuildID=dev-test'" .
