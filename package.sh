#!/usr/bin/env bash

rm -rf ./bin
GOOS=linux   GOARCH=amd64 go build -ldflags '-w -s' -o bin/tale-cli
upx bin/tale

#go build -o tale