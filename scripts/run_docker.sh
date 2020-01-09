#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)

docker run --name xgo-up -d xgo-up:latest
