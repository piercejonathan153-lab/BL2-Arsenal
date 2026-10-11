#!/bin/sh
# Builds dist/BL2ArsenalLauncher.exe (Windows x64) from any OS with Go installed.
set -e
cd "$(dirname "$0")"
go test ./...
go run github.com/tc-hib/go-winres@latest make --in winres/winres.json --arch amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -H windowsgui" -o dist/BL2ArsenalLauncher.exe .
echo built dist/BL2ArsenalLauncher.exe
