#!/bin/sh

filepath=$(
    cd $(dirname $0)
    pwd
)

docker exec -it xgo-up /bin/bash
