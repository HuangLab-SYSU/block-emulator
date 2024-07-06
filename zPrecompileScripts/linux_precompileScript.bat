@echo off
set GOOS=linux
set GOARCH=amd64
go build -o ../blockEmulator_Linux_Precompile ../main.go