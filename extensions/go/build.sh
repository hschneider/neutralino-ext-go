#!/bin/bash

DST=go

echo "Building $DST ..."
go build -ldflags="-s -w" -o $DST
chmod +x $DST
echo "DONE!"
