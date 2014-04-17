@echo off
@echo install packages...

set OLDGOPATH=%GOPATH%
set GOPATH=%cd%

go get -u github.com/robfig/revel
go get -u github.com/robfig/revel/revel
go get -u github.com/go-sql-driver/mysql
go get -u github.com/lunny/xorm
go get -u code.google.com/p/go-uuid/uuid
go get -u github.com/disintegration/imaging
go get -u github.com/cbonello/revel-csrf
go get -u github.com/qiniu/api
go get -u github.com/garyburd/redigo/redis
go get -u github.com/robfig/go-cache
go get -u github.com/robfig/gomemcache/memcache

set GOPATH=%OLDGOPATH%

echo finished.