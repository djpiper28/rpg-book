package pb_common

//go:generate echo Generating common protos...

//go:generate protoc --proto_path=../../protos/ --go_out=. --go_opt=paths=source_relative ../../protos/common.proto
