$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -ldflags="-X 'main.BuildID=dev-test'" -o smallkms.exe . && ./smallkms.exe admin "localhost:9001"
