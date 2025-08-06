#!/bin/bash

set -e

if [[ "$(uname)" == "Darwin" ]] || [[ "$(uname)" == "Linux" ]]; then
  bash -c "cd ./src/lib/grpcClient/pb/ && bash clean.sh"
  bash -c "protoc --ts_out ./src/lib/grpcClient/pb/ --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts --proto_path=../protos/ ../protos/*.proto"
else
  # We are on windows
  bash -c "cd ./src/lib/grpcClient/pb/ && bash clean.sh"
  bash -c "protoc --ts_out ./src/lib/grpcClient/pb/ --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts.cmd --proto_path=../protos/ ../protos/*.proto"
fi
