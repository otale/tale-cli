#!/usr/bin/bash

rm -rf bin
mkdir -p bin/linux_64
mkdir -p bin/macOSX_64

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o bin/linux_64/tale-cli && upx bin/linux_64/tale-cli
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o bin/macOSX_64/tale-cli && upx bin/macOSX_64/tale-cli
