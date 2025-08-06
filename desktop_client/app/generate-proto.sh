#!/bin/bash

args="--ts_out ./src/lib/grpcClient/pb/ --proto_path=../protos/"
proto_files="../protos/*.proto"

if [[ "$(uname)" == "Darwin" ]] || [[ "$(uname)" == "Linux" ]]; then
  protoc $args --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts $proto_files || exit 1
else
  protoc $args --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts $proto_files || exit 1
  # We are on windows
  # protoc $args --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts.cmd $proto_files || exit 1
fi
