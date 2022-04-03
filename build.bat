set GOOS=windows
set GOARCH=386

go install github.com/gonutz/rsrc@latest
rsrc -ico icon.ico

go build -ldflags="-H=windowsgui -s -w" -o "Schmetris.exe"
