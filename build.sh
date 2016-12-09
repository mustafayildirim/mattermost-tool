#!/bin/sh
CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o mattermosttool .
chmod a+x mattermosttool
