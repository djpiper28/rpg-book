package app

//go:generate bash -c "cd ./src/lib/grpcClient/pb/ && bash clean.sh"
//go:generate bash -c "protoc --ts_out ./src/lib/grpcClient/pb/ --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts --proto_path=../protos/ ../protos/*.proto"
