.PHONY: hooks fmt lint test build build-frontend install-tools

hooks:
	git config core.hooksPath .githooks

fmt:
	gofmt -w ./...

lint:
	golangci-lint run ./...

test:
	go test ./...

build:
	go build -o bin/zdzira ./cmd/zdzira

build-frontend:
	cd frontend && npm ci && npm run build

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
