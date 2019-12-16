#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)

export GO111MODULE=on
export GOPROXY=https://goproxy.io
export GOPATH=/home/songguang/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin
export GOROOT_BOOTSTRAP=$GOROOT

# cnpm install -g vue-cli
# cnpm install -g webpack
# cnpm install -g @vue/cli-init
# vue init webpack vue-demo

# ${GOROOT}/bin/go get github.com/zserge/webview
# ${GOROOT}/bin/go get github.com/spf13/viper

# for xgo
# sudo docker pull karalabe/xgo-1.12
# go get github.com/karalabe/xgo

docker build -f ${filepath}/xgo-update.dockfile -t xgo-update ${filepath}
# docker build -f ${filepath}/xgo-up.dockfile -t xgo-up ${filepath}
