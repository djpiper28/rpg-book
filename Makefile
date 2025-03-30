all: desktop_client

desktop-client-deps:
	cd ./desktop_client/app/ && pnpm i

desktop-client: desktop-client-deps go-generate
	cd ./desktop_client/app/ && pnpm build

go-generate:
	go generate ./...

go-core: go-generate

go-test: go-core
	go test ./...

desktop_client: go-core desktop-client
	cd desktop_client/launcher && go build
