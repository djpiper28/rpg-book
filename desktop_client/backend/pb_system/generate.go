package pb_system

//go:generate echo Generating system protos...

//go:generate protoc --proto_path=../../protos/ --go_out=. --go_opt=paths=source_relative ../../protos/system.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
