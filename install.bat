@echo off
@echo install packages...

set OLDGOPATH=%GOPATH%
set GOPATH=%cd%

go get -u github.com/robfig/revel
go get -u github.com/robfig/revel/revel
go get -u github.com/coocood/qbs
go get -u github.com/coocood/mysql
go get -u code.google.com/p/go-uuid/uuid
go get -u github.com/disintegration/imaging
go get -u github.com/cbonello/revel-csrf

set GOPATH=%OLDGOPATH%

echo finished.