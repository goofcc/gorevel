#!/usr/bin/env bash

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

bin/revel run gorevel
# 可选参数：[dev]/[prod] + [port]
# bin/revel run gorevel dev 3000
# bin/revel run gorevel prod 3000

export GOPATH="$OLDGOPATH"
