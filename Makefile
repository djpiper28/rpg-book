.PHONY: all
all: desktop-client

DESKTOP_APP=$(shell pwd)/desktop_client/app/

.PHONY: desktop-client-deps
desktop-client-deps:
	cd $(DESKTOP_APP) && pnpm i

.PHONY: desktop-client-codegen
desktop-client-codegen: go-generate

.PHONY: desktop-client-app
desktop-client-app: desktop-client-codegen 
	cd $(DESKTOP_APP) && pnpm build

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

.PHONY: test
test: go-test desktop-client-test

.PHONY: go-fmt
go-fmt:
	go fmt -w -l ./...

.PHONY: prettier
prettier:
	cd $(DESKTOP_APP) && npx prettier -w .

.PHONY: format
format: go-fmt prettier

.PHONY: cleanup
cleanup:
	find desktop_client | grep -E ".*\\.sqlite(-journal)?" | xargs rm -rf
