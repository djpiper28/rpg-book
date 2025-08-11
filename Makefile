.PHONY: all
all: desktop

DESKTOP_APP="$(shell pwd)/desktop_client/app/"

.PHONY: desktop-deps
desktop-deps:
	cd $(DESKTOP_APP) && pnpm i

.PHONY: desktop-codegen
desktop-codegen: go-generate

.PHONY: desktop-app
desktop-app: desktop-codegen 
	cd $(DESKTOP_APP) && pnpm build

# Starts a dev electron app (pnpm calls other targets in this, this just makes it simpler to use)
.PHONY: dev
dev:
	cd $(DESKTOP_APP) && pnpm dev

# This only works on linux
.PHONY: release
release: all
	./desktop_client/launcher/launcher ./desktop_client/app/release/1.0.0/RPG-Book-Linux-1.0.0.AppImage

.PHONY: desktop-backend
desktop-backend: go-core
	cd ./desktop_client/launcher && go build

.PHONY: go-generate
go-generate: desktop-deps
	go generate ./...

.PHONY: go-core
go-core: go-generate

.PHONY: go-test
go-test: go-core
	go test ./...

.PHONY: desktop
desktop: go-core desktop-app desktop-backend

.PHONY: desktop-test
desktop-test: desktop-codegen 
	cd $(DESKTOP_APP) && pnpm test

.PHONY: go-lint
go-lint: go-core
	go vet ./...

.PHONY: test
test: go-test go-lint desktop-test

.PHONY: go-fmt
go-fmt:
	gofmt -w -l .

.PHONY: prettier
prettier:
	cd $(DESKTOP_APP) && pnpm lint --fix
	cd $(DESKTOP_APP) && npx prettier -w .

.PHONY: format
format: go-fmt prettier

.PHONY: cleanup
cleanup:
	find desktop_client | grep -E ".*\\.sqlite(-journal)?" | xargs rm -rf
