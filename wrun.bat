
@echo off

set OLDGOPATH=%GOPATH%
set GOPATH=%cd%

bin\revel run gorevel
:: 可选参数：[dev]/[prod] + [port]
:: bin/revel run gorevel dev 9000
:: bin/revel run gorevel prod 9000

set GOPATH=%OLDGOPATH%
