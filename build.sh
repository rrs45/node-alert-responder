#!/usr/bin/env bash

set -ex

godir=/tmp/go/src/github.com

mkdir -p $godir/box-node-alert-responder
export GOPATH=/tmp/go

mv cmd $godir/box-node-alert-responder/
mv pkg $godir/box-node-alert-responder/
cd $godir/box-node-alert-responder
mkdir bin
go get ./...
CGO_ENABLED=0 GOOS=linux go build -o bin/node-alert-responder -ldflags '-w' cmd/node-alert-responder.go

mkdir -p /git-root/build
mv bin/node-alert-responder /git-root/build