.PHONY: lint test generate producer consumer

-include .env
export

lint:
	@golangci-lint run

test:
	@go test ./...

generate:
	@go generate ./...

