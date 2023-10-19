$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -ldflags="-X 'main.BuildID=dev-test'" -o agent.exe . && ./agent.exe -env .env server :10443
