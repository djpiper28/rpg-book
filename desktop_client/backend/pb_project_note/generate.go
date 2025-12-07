package pb_project_note

//go:generate echo Generating project note protos...

//go:generate protoc --proto_path=../../protos/ --go_out=. --go_opt=paths=source_relative ../../protos/project_note.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative
