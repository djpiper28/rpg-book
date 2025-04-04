package pb_settings

//go:generate protoc --proto_path=../../protos/ --go_out=. --go_opt=paths=source_relative ../../protos/settings.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
