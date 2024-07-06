#!/bin/bash

GOOS=darwin \
GOARCH=amd64 \
go build -o ../blockEmulator_MacOS_Precompile ../main.go
