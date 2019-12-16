#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)

# cnpm install -g vue-cli
# cnpm install -g webpack
# cnpm install -g @vue/cli-init
# vue init webpack vue-demo

# ${GOROOT}/bin/go get github.com/zserge/webview
# ${GOROOT}/bin/go get github.com/spf13/viper

# for xgo
# sudo docker pull karalabe/xgo-1.12
# go get github.com/karalabe/xgo

# docker build -f ${filepath}/xgo-update.dockfile -t xgo-update ${filepath}
# docker build -f ${filepath}/xgo-up.dockfile -t xgo-up ${filepath}

# docker run --name xgo-update -d xgo-update
docker exec -it xgo-update /bin/bash
