@echo off 



set CURPATH=%cd%

set OLDGOPATH=%GOPATH%
set GOPATH=%CURPATH%;%OLDGOPATH%

echo %GOPATH%
go install publicIP

@echo on 