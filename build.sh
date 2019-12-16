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

if [ -z $1 ]; then
  compileType="linux"
else
  compileType=$1
fi
echo $compileType
# exit 1

# GOPATH=/home/songguang/go
# GOROOT=/usr/local/go

# cd ${GOPATH}/src/ssd_helper_webview/

export GO111MODULE=on
export GOPROXY=https://goproxy.io
export GOPATH=/home/songguang/go
export GOROOT=/usr/local/go
# export PATH=$PATH:$GOPATH/bin
# export GOROOT_BOOTSTRAP=$GOROOT

# ${GOROOT}/bin/go env

# rm -rf ${filepath}/app/*

if [ ! -d "${GOPATH}/src/ssd_helper_webview" ]; then
  mkdir ${GOPATH}/src/ssd_helper_webview
fi
rm -rf ${GOPATH}/src/ssd_helper_webview/*
cp -R $filepath/* ${GOPATH}/src/ssd_helper_webview/
cp -R ${filepath}/static/vue-demo/dist/* ${GOPATH}/src/ssd_helper_webview/app/

if [ $compileType = "cross" ]; then
  cd ${GOPATH}/src/ssd_helper_webview/ && ${GOPATH}/bin/xgo --go=1.12 --targets="linux/amd64,windows-7.0/amd64" --image=xgo-update ssd_helper_webview
  # cd ${GOPATH}/src/ssd_helper_webview/ && ${GOPATH}/bin/xgo --targets="linux/amd64,windows-7.0/amd64" ./
else
  cd ${GOPATH}/src/ssd_helper_webview/ && ${GOROOT}/bin/go build -o ssd_helper_webview main.go
fi

rm -rf $filepath/dist/*

cp ${GOPATH}/src/ssd_helper_webview/ssd_helper_webview* $filepath/dist/
cp -R ${GOPATH}/src/ssd_helper_webview/app $filepath/dist/
cp -R ${GOPATH}/src/ssd_helper_webview/conf $filepath/dist/conf
