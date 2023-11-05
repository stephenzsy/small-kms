$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -ldflags="-X 'main.BuildID=dev-test'" -o agent.exe . && ./agent.exe --pretty-log -env .env server :10443
