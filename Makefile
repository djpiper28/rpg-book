all: desktop_client

go-generate:
	go generate ./...

go-core: go-generate

go-test: go-core
	go test ./...

desktop_client: go-core
	cd desktop_client/launcher && go build
