#!/bin/bash

GOOS=linux \
GOARCH=amd64 \
go build -o ../blockEmulator_Linux_Precompile ../main.go
