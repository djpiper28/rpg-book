.PHONY: all
all: desktop-client

DESKTOP_APP="$(shell pwd)/desktop_client/app/"

.PHONY: desktop-client-deps
desktop-client-deps:
	cd $(DESKTOP_APP) && pnpm i

.PHONY: desktop-client-codegen
desktop-client-codegen: go-generate

.PHONY: desktop-client-app
desktop-client-app: desktop-client-codegen 
	cd $(DESKTOP_APP) && pnpm build

# Starts a dev electron app (pnpm calls other targets in this, this just makes it simpler to use)
.PHONY: dev
dev:
	cd $(DESKTOP_APP) && pnpm dev

.PHONY: desktop-client-backend
desktop-client-backend: go-core
	cd ./desktop_client/launcher && go build

.PHONY: go-generate
go-generate: desktop-client-deps
	PATH="$$PATH:$(DESKTOP_APP)/node_modules/.bin" go generate ./...

.PHONY: go-core
go-core: go-generate

.PHONY: go-test
go-test: go-core
	go test ./...

.PHONY: desktop-client
desktop-client: go-core desktop-client-app desktop-client-backend

.PHONY: desktop-client-test
desktop-client-test: desktop-client-codegen 
	cd $(DESKTOP_APP) && pnpm test

.PHONY: go-lint
go-lint: go-core
	go vet ./...

.PHONY: test
test: go-test go-lint desktop-client-test

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
