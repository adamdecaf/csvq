ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

.PHONY: build
build:
	go build -o ./bin/csvq -ldflags "-X github.com/adamdecaf/csvq.Version=${VERSION}" github.com/adamdecaf/csvq/cmd/csvq

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=0.0 ./lint-project.sh
endif

build-webui:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js ./docs/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o ./docs/csvq.wasm github.com/adamdecaf/csvq/internal/webui/

.PHONY: clean
clean:
	@rm -rf ./bin/ ./tmp/ coverage.txt misspell* staticcheck lint-project.sh

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out
