#!/usr/bin/env bash
echo 'install packages...'

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

go get -u github.com/robfig/revel
go get -u github.com/robfig/revel/revel
go get -u github.com/coocood/qbs
go get -u github.com/coocood/mysql
go get -u code.google.com/p/go-uuid/uuid
go get -u github.com/disintegration/imaging
go get -u github.com/cbonello/revel-csrf

export GOPATH="$OLDGOPATH"

echo 'finished.'