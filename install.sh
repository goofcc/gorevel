#!/usr/bin/env bash
echo 'install packages...'

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

go get -u -v github.com/robfig/revel
go get -u -v github.com/robfig/revel/revel
go get -u -v github.com/go-sql-driver/mysql
go get -u -v github.com/go-xorm/xorm
go get -u -v github.com/pborman/uuid
go get -u -v github.com/disintegration/imaging
go get -u -v github.com/cbonello/revel-csrf
go get -u -v github.com/qiniu/api
go get -u -v github.com/garyburd/redigo/redis
go get -u -v github.com/robfig/go-cache
go get -u -v github.com/robfig/gomemcache/memcache

export GOPATH="$OLDGOPATH"

echo 'finished.'
