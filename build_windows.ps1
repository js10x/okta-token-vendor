$bin_name   = "oktv.exe"
$Env:GOOS   = "windows"
$Env:GOARCH = "amd64"
go build -ldflags=-w -o $bin_name "main.go"
