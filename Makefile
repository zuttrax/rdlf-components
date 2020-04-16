export GO111MODULE ?= on
PACKAGES = $(shell go list ./...)

.PHONY: all
all: deps test fmt linter

.PHONY: deps
deps:
	@echo "Ensuring deps"
	@go mod tidy

.PHONY: test
test:
	@echo "Running tests"
	@go test ./... -covermode=atomic

.PHONY: fmt
fmt:
	@echo "Executing fmt"
	@go fmt $(PACKAGES)

.PHONY: linter
linter:
	@echo "Executing golangci-lint"
	@golangci-lint run