$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -o smallkms.exe . && ./smallkms.exe -env .env --pretty-log admin "localhost:9001"
