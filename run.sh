#!/usr/bin/env bash

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

bin/revel run gorevel
# bin/revel run gorevel dev 9000
# bin/revel run gorevel prod 9000

export GOPATH="$OLDGOPATH"
