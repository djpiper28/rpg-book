.PHONY: all
all: desktop-client

.PHONY: desktop-client-deps
desktop-client-deps:
	cd ./desktop_client/app/ && pnpm i

.PHONY: desktop-client-app
desktop-client-app: desktop-client-deps go-generate
	cd ./desktop_client/app/ && pnpm build

.PHONY: desktop-client-backend
desktop-client-backend: go-core
	cd ./desktop_client/launcher && go build

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: go-core
go-core: go-generate

.PHONY: go-test
go-test: go-core
	go test ./...

.PHONY: desktop-client
desktop-client: go-core desktop-client-app desktop-client-backend
