package pb_project

//go:generate echo Generating project protos...

//go:generate protoc --proto_path=../../protos/ --go_out=. --go_opt=paths=source_relative ../../protos/project.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
