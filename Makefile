all: desktop_client

desktop-client-deps:
	cd ./desktop_client/app/ && pnpm i

desktop-client-app: desktop-client-deps go-generate
	cd ./desktop_client/app/ && pnpm build

desktop-client-backend: go-core
	cd ./desktop_client/launcher && go build

go-generate:
	go generate ./...

go-core: go-generate

go-test: go-core
	go test ./...

desktop-client: go-core desktop-client-app desktop-client-backend
