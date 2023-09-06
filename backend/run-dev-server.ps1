go build -o smallkms.exe .
get-content .env | foreach-object {
    $name, $value = $_.split('=')
    if (![string]::IsNullOrWhiteSpace($name) && !$name.Contains('#')) {
        Set-Content env:\$name $value
    }
}
./smallkms.exe admin "localhost:9001"
