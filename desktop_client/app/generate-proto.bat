@echo off
set args=--ts_out ./src/lib/grpcClient/pb/ --proto_path=../protos/
set proto_files=../protos/*.proto

protoc %args% --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts.cmd %proto_files% || exit 1
