#!/usr/bin/env bash

mkdir -p ./bin
export GIN_MODE=release
cp -rf "./config.toml" "./bin/config.toml"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o ./bin/rpc_demo .