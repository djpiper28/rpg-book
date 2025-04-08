package app

//go:generate echo Generating JS protos...

//go:generate bash -c "cd ./src/lib/grpcClient/pb/ && bash clean.sh"

//go:generate bash -c "protoc --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin --js_out=import_style=commonjs_strict,binary:./src/lib/grpcClient/pb/ --ts_out=service=grpc-node,mode=grpc-js:./src/lib/grpcClient/pb/ --grpc_out=grpc_js:./src/lib/grpcClient/pb/ --proto_path=../protos/ ../protos/*.proto"
