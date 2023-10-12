$env:GOARCH="amd64"
$env:GOOS="windows"
go build -o smallkms.exe . &&  ./smallkms.exe admin "localhost:9001"
