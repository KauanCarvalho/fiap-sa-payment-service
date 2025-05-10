include .env

APP_NAME := fiap_sa_payment_service
BDD_TEST_DIR=internal/test/bdd
BIN_DIR := bin
DOCKER_COMPOSE := docker-compose
ENV_FILE := .env
GO ?= go

.DEFAULT_GOAL := help
.PHONY: help deps setup-git-hooks lint check-coverage test test-bdd test-payment-service coverage-html build-api run-api run-api-air docker-up docker-down swag build-worker run-worker run-worker-air

help:
	@echo ""
	@echo "Available targets:"
	@echo "  help                  # Show this help message"
	@echo "  deps                  # Install dependencies"
	@echo "  setup-git-hooks       # Install Git hooks using Lefthook"
	@echo "  lint                  # Run linters"
	@echo "  check-coverage        # Check test coverage"
	@echo "  test                  # Run tests"
	@echo "  test-bdd              # Run BDD tests"
	@echo "  test-payment-service  # Run tests for payment service"
	@echo "  coverage-html         # Generate HTML coverage report"
	@echo "  build-api             # Build the API"
	@echo "  run-api               # Run the API"
	@echo "  run-api-air           # Run the API with live reloading"
	@echo "  build-worker          # Build the worker"
	@echo "  run-worker            # Run the worker"
	@echo "  run-worker-air        # Run the worker with live reloading"
	@echo "  docker-up             # Start Docker container(s)"
	@echo "  docker-down           # Stop Docker containers"
	@echo "  swag                  # Generate Swagger documentation"
	@echo ""

deps:
	@echo "Installing dependencies..."
	$(GO) mod download

setup-git-hooks: deps
	@echo "Installing Git hooks with Lefthook..."
	$(GO) tool lefthook install

lint: deps
	@echo "Running linter..."
	$(GO) tool golangci-lint run ./... --config .golangci.yml

test:
	@echo "Running tests..."
	DB_NAME=$(DB_NAME)_test $(GO) tool godotenv -f $(ENV_FILE) $(GO) test ./... -coverprofile=coverage.out -cover -p 1

test-bdd:
	@echo "Running BDD tests..."
	DB_NAME=$(DB_NAME)_test APP_ENV=test $(GO) tool godotenv -f $(ENV_FILE) $(GO) test -v ./$(BDD_TEST_DIR) --tags=integration --coverprofile=coverage.out --cover -p 1

test-payment-service:
	@echo "Running tests for payment service..."
	@./testdata/test-payment-service.sh $(filter-out $@,$(MAKECMDGOALS))

check-coverage: test
	@echo "Checking coverage..."
	$(GO) tool go-test-coverage --config=./.testcoverage.yml

coverage-html: test
	@echo "Openning coverage report..."
	$(GO) tool cover -html=coverage.out

build-api:
	@echo "Building api..."
	$(GO) build -o $(BIN_DIR)/$(APP_NAME)_api ./cmd/api/main.go

run-api: build-api
	@echo "Running api..."
	$(BIN_DIR)/$(APP_NAME)_api

run-api-air: deps
	@echo "Running api with live reloading..."
	$(GO) tool air -c .air.api.toml

build-worker:
	@echo "Building worker..."
	$(GO) build -o $(BIN_DIR)/$(APP_NAME)_worker ./cmd/worker/main.go

run-worker: build-worker
	@echo "Running worker..."
	$(BIN_DIR)/$(APP_NAME)_worker

run-worker-air: deps
	@echo "Running worker with live reloading..."
	$(GO) tool air -c .air.worker.toml

docker-up:
	@echo "Starting Docker container(s)..."
	$(DOCKER_COMPOSE) up -d $(filter-out $@,$(MAKECMDGOALS))

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

swag:
	@echo "Generating Swagger documentation..."
	$(GO) tool swag init --parseDependency --parseInternal -g cmd/api/main.go -o ./swagger --ot json,go
