#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)
parentpath=$(
    cd $(dirname $filepath)
    pwd
)
rootpath=$(
    cd $(dirname $parentpath)
    pwd
)

export GO111MODULE=on
export GOPROXY=https://goproxy.io
export GOPATH=/home/songguang/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin
export GOROOT_BOOTSTRAP=$GOROOT

# ${GOROOT}/bin/go env

${GOROOT}/bin/go run .
