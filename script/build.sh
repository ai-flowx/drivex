#!/bin/bash

ldflags="-s -w"
target="drivex"

go env -w GOPROXY=https://goproxy.io,direct

CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target main.go
