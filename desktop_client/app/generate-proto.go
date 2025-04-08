package app

//go:generate bash -c "cd ./src/lib/grpcClient/pb/ && bash clean.sh"
//go:generate bash -c "protoc --js_out=import_style=es6,binary:./src/lib/grpcClient/pb/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/lib/grpcClient/pb/ --proto_path=../protos/ ../protos/*.proto"

//--js_out: Unknown import style poopoo, expected one of: closure, commonjs, browser, es6.
