@echo off

set OLDGOPATH=%GOPATH%
set GOPATH=%cd%

:: 可选参数：生产模式[prod]，端口[8080]
bin\revel run gorevel

set GOPATH=%OLDGOPATH%
