BINARY    := calendar
MAIN_FILE := cmd/$(BINARY)/main.go

PKG       := $(shell go list -m)
VERSION   := $(or $(IMAGE_TAG),$(shell git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0))

LOCAL_BIN_DIR := $(PWD)/bin

DOCKER_COMPOSE := docker compose --project-name=$(BINARY) --file=env/docker-compose.yaml

.DEFAULT_GOAL := help

.PHONY: run
run: ## Run the application
	go run $(MAIN_FILE)

.PHONY: env
env: ## Create environment
	@echo "> Creating environment $(BINARY)"
	$(DOCKER_COMPOSE) up -d

.PHONY: env-down
env-down: ## Destroy environment
	@echo "> Destroying environment $(BINARY)"
	$(DOCKER_COMPOSE) down --volumes

# go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" cmd/$(BINARY)/main.go
.PHONY: build
build: ## Build the binary
	goreleaser build --snapshot --clean --single-target

.PHONY: docs
docs: ## Generate Swagger documentation
	go mod download -x
	swag init -pd -g internal/server/server.go -o internal/server/docs

.PHONY: lint
lint: ## Lint Go files
	@GOPATH="$(shell dirname $(PWD))" golangci-lint run ./...

.PHONY: test
test: ## Run unit tests
	@go test -v -race ./...

.PHONY: coverage
coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
