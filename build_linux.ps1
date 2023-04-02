$bin_name   = "oktv.bin"
$Env:GOOS   = "linux"
$Env:GOARCH = "amd64"
go build -ldflags=-w -o $bin_name "main.go"