@echo off

set DST=go.exe

echo Building %DST% ...
go build -ldflags="-s -w" -o %DST%

echo DONE!
