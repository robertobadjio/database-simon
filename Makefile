#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

export GOOS=linux
export GOARCH=amd64

GO_TEST_COMMAND = go test


mockgen: ## Run mockgen
	# go install go.uber.org/mock/mockgen@latest
	mockgen -source=./internal/database/database.go -destination=./internal/database/database_mock.go -package=database
	mockgen -source=./internal/database/compute/compute.go -destination=./internal/database/compute/compute_mock.go -package=compute
	mockgen -source=./internal/database/storage/storage.go -destination=./internal/database/storage/storage_mock.go -package=storage
	mockgen -source=./internal/database/storage/engine.go -destination=./internal/database/storage/engine_mock.go -package=storage

test-unit: ## Run unit tests
	$(GO_TEST_COMMAND) \
		./internal/...
		-count=1 \
		-cover -coverprofile=$(TEST_COVER_FILENAME)

test-unit-race: ## Run unit tests with -race flag
	$(GO_TEST_COMMAND) ./internal/... -count=1 -race