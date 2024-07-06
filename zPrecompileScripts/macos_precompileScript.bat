@echo off
set GOOS=darwin
set GOARCH=amd64
go build -o ../blockEmulator_MacOS_Precompile ../main.go