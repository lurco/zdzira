.PHONY: hooks fmt lint test build build-frontend install-tools openapi

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

# Regenerate the committed OpenAPI snapshot from the live route definitions.
openapi:
	go run ./cmd/zdzira -db ":memory:" -dump-openapi > docs/openapi.json
