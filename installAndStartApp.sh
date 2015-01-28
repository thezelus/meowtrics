#!/bin/bash

go get github.com/tools/godep 
cd model
godep restore
go build
go install
cd ..
cd server
godep restore
go build
bash startServer.sh
