#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

export GOOS=linux
export GOARCH=amd64

GO_TEST_COMMAND = go test
TEST_COVER_FILENAME = c.out

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps-linux: ## Install dependencies for Linux
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v1.64.5

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt lint-config-verify ## Check code (lint)
	./bin/golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

lint-config-verify: fmt ## Verify config (lint)
	./bin/golangci-lint config verify --config .golangci.pipeline.yaml

vet: fmt ## Check code (vet)
	go vet ./...
.PHONY:vet

vet-shadow: fmt ## Check code with detect shadow (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

mockgen: ## Run mockgen
	# go install go.uber.org/mock/mockgen@latest
	mockgen -source=./internal/database/database.go -destination=./internal/database/database_mock.go -package=database
	mockgen -source=./internal/database/compute/compute.go -destination=./internal/database/compute/compute_mock.go -package=compute
	mockgen -source=./internal/database/storage/storage.go -destination=./internal/database/storage/storage_mock.go -package=storage
	mockgen -source=./internal/database/storage/engine.go -destination=./internal/database/storage/engine_mock.go -package=storage
	mockgen -source=./internal/database/storage/engine.go -destination=./internal/database/storage/engine_mock.go -package=storage
	mockgen -source=./internal/config/enviroment.go -destination=./internal/config/enviroment_mock.go -package=config

test-unit: ## Run unit tests
	$(GO_TEST_COMMAND) \
		./internal/... \
		-count=1 \
		-cover -coverprofile=$(TEST_COVER_FILENAME)

test-unit-race: ## Run unit tests with -race flag
	$(GO_TEST_COMMAND) ./internal/... -count=1 -race