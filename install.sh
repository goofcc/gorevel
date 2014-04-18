#!/usr/bin/env bash
echo 'install packages...'

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

go get -u github.com/robfig/revel
go get -u github.com/robfig/revel/revel
go get -u github.com/go-sql-driver/mysql
go get -u github.com/go-xorm/xorm
go get -u code.google.com/p/go-uuid/uuid
go get -u github.com/disintegration/imaging
go get -u github.com/cbonello/revel-csrf
go get -u github.com/qiniu/api
go get -u github.com/garyburd/redigo/redis
go get -u github.com/robfig/go-cache
go get -u github.com/robfig/gomemcache/memcache

export GOPATH="$OLDGOPATH"

echo 'finished.'