#!/bin/bash

GOOS=darwin \
GOARCH=amd64 \
go build -o ../blockEmulator_darwin_Precompile ../main.go
