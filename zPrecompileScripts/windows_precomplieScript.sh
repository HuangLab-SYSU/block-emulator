#!/bin/bash

GOOS=windows \
GOARCH=amd64 \
go build -o ../blockEmulator_Windows_Precompile.exe ../main.go
