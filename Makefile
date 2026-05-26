.PHONY: hooks fmt lint test build install-tools

hooks:
	git config core.hooksPath .githooks

fmt:
	gofmt -w ./...

lint:
	golangci-lint run ./...

test:
	go test ./...

build:
	go build -o bin/issuetracker ./cmd/issuetracker

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
