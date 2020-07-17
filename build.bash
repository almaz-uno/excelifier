#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o excelifier.linux
GOOS=darwin GOARCH=amd64 go build -o excelifier.darwin
GOOS=windows GOARCH=amd64 go build -o excelifier.windows.exe