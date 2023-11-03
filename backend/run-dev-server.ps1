$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -o smallkms.exe . && ./smallkms.exe -env .env --debug --pretty-log admin "localhost:9001"
