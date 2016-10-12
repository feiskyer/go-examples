#!/bin/sh
# git clone https://github.com/google/protobuf
# cd protobuf
# ./autogen.sh
# ./configure
# make
# make install

# go get -u google.golang.org/grpc
# go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:. logs.proto
