export GO111MODULE ?= on

.PHONY: all
all: test linter

.PHONY: test
test:
	@echo "Running tests"
	@go test ./... -covermode=atomic

.PHONY: linter
linter:
	@echo "Executing golangci-lint"
	@golangci-lint run