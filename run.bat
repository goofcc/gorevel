@echo off

set OLDGOPATH=%GOPATH%
set GOPATH=%cd%

bin\revel run gorevel

set GOPATH=%OLDGOPATH%
