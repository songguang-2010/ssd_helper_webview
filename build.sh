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
#默认linux平台编译
  compileType="linux"
else
#根据参数编译目标平台
  compileType=$1
fi
echo $compileType
# exit 1

# GOPATH=/home/songguang/go
# GOROOT=/usr/local/go

# cd ${GOPATH}/src/ssd_helper_webview/

#启用go的module编译模式
export GO111MODULE=on
#go代理
export GOPROXY=https://goproxy.io
export GOPATH=/home/songguang/go
export GOROOT=/usr/local/go
# export XGO_IN_XGO=1
# export PATH=$PATH:$GOPATH/bin
# export GOROOT_BOOTSTRAP=$GOROOT

# ${GOROOT}/bin/go env

# rm -rf ${filepath}/app/*

#补全编译目录及清除旧数据
if [ ! -d "${GOPATH}/src/ssd_helper_webview" ]; then
  mkdir ${GOPATH}/src/ssd_helper_webview
fi
rm -rf ${GOPATH}/src/ssd_helper_webview/*
#拷贝相关文件至编译目录
cp -R $filepath/* ${GOPATH}/src/ssd_helper_webview/
cp -R ${filepath}/static/vue-demo/dist/* ${GOPATH}/src/ssd_helper_webview/app/

#交叉编译
if [ $compileType = "cross" ]; then
  cd ${GOPATH}/src/ && ${GOPATH}/bin/xgo --image=xgo-up --targets="linux/amd64,windows-7.0/amd64" -dest ./ssd_helper_webview/ ./ssd_helper_webview
  # cd ${GOPATH}/src/ && ${GOPATH}/bin/xgo --image=xgo-up --targets="linux/amd64,windows-7.0/amd64" -dest ./ssd_helper_webview/ ./ssd_helper_webview
else
#linux编译
  cd ${GOPATH}/src/ssd_helper_webview/ && ${GOROOT}/bin/go build -o ssd_helper_webview main.go
fi

rm -rf $filepath/dist/*

#编译结果拷贝至目标输出目录
cp ${GOPATH}/src/ssd_helper_webview/ssd_helper_webview* $filepath/dist/
cp -R ${GOPATH}/src/ssd_helper_webview/app $filepath/dist/
cp -R ${GOPATH}/src/ssd_helper_webview/conf $filepath/dist/conf
