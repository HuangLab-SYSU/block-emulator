#!/bin/bash

GOOS=linux \
GOARCH=amd64 \
go build -o ../blockEmulator_linux_Precompile ../main.go
