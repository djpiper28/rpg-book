.PHONY: all
all: desktop

DESKTOP_APP="$(shell pwd)/desktop_client/app/"
GOEXPERIMENT=nodwarf5,jsonv2

.PHONY: desktop-deps
desktop-deps:
	cd $(DESKTOP_APP) && pnpm i

.PHONY: desktop-codegen
desktop-codegen: go-generate

.PHONY: desktop-app
desktop-app: desktop-codegen desktop-backend 
	cd $(DESKTOP_APP) && pnpm build

# Starts a dev electron app (pnpm calls other targets in this, this just makes it simpler to use)
.PHONY: dev
dev: desktop-deps desktop-backend
	cd $(DESKTOP_APP) && pnpm dev

# This only works on linux
.PHONY: release
release: all
	./desktop_client/app/release/1.0.0/RPG-Book-Linux-1.0.0.AppImage

.PHONY: desktop-backend
desktop-backend: go-core
	cd ./desktop_client/launcher && go build

.PHONY: go-generate
go-generate: desktop-deps
	go generate ./...

.PHONY: go-core
go-core: go-generate

.PHONY: fuzz
fuzz: go-core
	go test -fuzz=FuzzParser ./common/search/parser/
	go test -fuzz=FuzzSearchCharacter,FuzzSearchNote ./desktop_client/backend/project/
	$(MAKE) clean

.PHONY: go-test
go-test: go-core
	go test ./...

.PHONY: desktop
desktop: desktop-app desktop-backend

.PHONY: desktop-test
desktop-test: desktop-codegen 
	cd $(DESKTOP_APP) && pnpm test

.PHONY: desktop_lint
desktop_lint: desktop-codegen 
	cd $(DESKTOP_APP) && pnpm lint

.PHONY: go-lint
go-lint: go-core
	go vet ./...

.PHONY: test
test: go-test go-lint desktop_lint desktop-test
	$(MAKE) clean

.PHONY: go-fmt
go-fmt:
	gofumpt -w -l .

.PHONY: prettier
prettier:
	cd $(DESKTOP_APP) && pnpm lint --fix
	cd $(DESKTOP_APP) && npx prettier -w .
	npx prettier -w *.md

.PHONY: format
format: go-fmt prettier

.PHONY: clean-node-modules
clean-node-modules:
	rm -rf ./desktop_client/app/node_modules/

.PHONY: clean
clean: clean-node-modules
	find . | grep -E ".*\\.(sqlite|rpg)((-journal)|(-wal)|(-shm))?" | xargs rm -rf -d '\n'
	find . | grep -E "testdata/.*/[a-z0-9]+" | xargs rm -rf -d '\n'
