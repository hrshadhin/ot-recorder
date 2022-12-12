#!/usr/bin/env bash

BUILD_VERSION := $(shell git describe --always --tags)
BUILD_TIME=$(shell date '+%Y%m%d-%H%M%S')
DOCKER_IMAGE_NAME="hrshadhin/ot-recoder"
BINARY_NAME=ot-recoder
BIN_OUT_DIR=bin

export PATH=$(shell go env GOPATH)/bin:$(shell echo $$PATH)

.PHONY: all

all: dl-deps build test-unit ## Build binary (with unit tests)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: build ## Run lint checks
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(shell go env GOPATH)/bin/golangci-lint run

fmt: ## Refactor go files with gofmt and goimports
	go install golang.org/x/tools/cmd/goimports@latest
	find . -name '*.go' | while read -r file; do goimports -w "$$file"; done

test-unit:  ## Run unit tests
	go test -v -coverprofile=coverage.txt -covermode=atomic -cover ./app/...

test-integration:  ## Run sqlite integration tests
	go test -v -tags=integration ./it -count=1

test-integration-mysql:  ## Run mysql integration tests
	go test -v -tags=integration ./it/mysql -count=1

test-integration-pgsql:  ## Run pgsql integration tests
	go test -v -tags=integration ./it/pgsql -count=1

clean: ## Cleans output directory
	$(shell rm -rf $(BIN_OUT_DIR)/*)
	$(shell rm -rf cmd/migrations)
	$(shell rm -rf ./*.db ./it/*.db coverage.txt _doc/docs.go _doc/swagger.json _doc/swagger.yaml)

dl-deps: ## Get dependencies
	go mod vendor

build: clean ## Build binary
	go generate ./cmd
	go build -v -ldflags="-w -s -X ot-recorder/app.Version=${BUILD_VERSION} -X ot-recorder/app.BuildTime=${BUILD_TIME}" -o $(BIN_OUT_DIR)/$(BINARY_NAME)

version: ## Check binary version
	./$(BIN_OUT_DIR)/$(BINARY_NAME) --version

serve: build ## Run http server
	./$(BIN_OUT_DIR)/$(BINARY_NAME) serve

doc: ## Creates swagger documentation as html file
	go install github.com/swaggo/swag/cmd/swag@v1.8.4
	$(shell go env GOPATH)/bin/swag init -g _doc/api.go -o _doc
	$(shell which redoc-cli) build --options.disableSearch -o _doc/swagger.html _doc/swagger.json

migrate-up: ## Run migration
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate up

migrate-down: ## Revert migration
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate down

docker-build: ## Build docker image
	docker build --build-arg BUILD_VERSION=${BUILD_VERSION} --build-arg BUILD_TIME=${BUILD_TIME} --tag ${DOCKER_IMAGE_NAME} .

docker-push: ## Push docker image
	docker push

docker-run: ## Run docker image with sqlite
	mkdir -p data
	sudo chown -R 1000:1000 data
	docker run --name ot-recoder --rm -it -p 8000:8000 \
		-v $$(pwd)/config.yml:/app/config.yml \
		-v $$(pwd)/data:/persist \
		$(DOCKER_IMAGE_NAME):latest

docker-migrate: ## Run migrations inside dorker
	docker exec ot-recoder /app/ot-recoder migrate up
