set GOOS=windows

go build -o shell_win.exe -ldflags="-s -w" ./cmd/client