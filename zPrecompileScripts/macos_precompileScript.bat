@echo off
set GOOS=darwin
set GOARCH=amd64
go build -o ../blockEmulator_darwin_Precompile ../main.go