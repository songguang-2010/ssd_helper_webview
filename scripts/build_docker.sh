#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)

docker build -f ${filepath}/xgo-up.dockfile -t xgo-up ${filepath}
