.PHONY: help build run test test-coverage lint fmt vet clean migrate-up migrate-down docker-up docker-down install-tools install-hooks

# Variables
APP_NAME := idiomapi
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin
MAIN_FILE := $(CMD_DIR)/main.go
BINARY := $(BUILD_DIR)/$(APP_NAME)
CONFIG_FILE := configs/config.toml

# Database variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= tododb
DB_DSN := "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable"

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

## help: Display this help message
help:
	@echo "$(CYAN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

## build: Build the application binary
build:
	@echo "$(CYAN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) $(MAIN_FILE)
	@echo "$(GREEN)Build complete: $(BINARY)$(NC)"

## run: Run the application
run:
	@echo "$(CYAN)Running $(APP_NAME)...$(NC)"
	@go run $(MAIN_FILE) -config $(CONFIG_FILE)

## clean: Remove build artifacts
clean:
	@echo "$(CYAN)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "$(GREEN)Clean complete$(NC)"

##@ Testing

## test: Run tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	@go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

##@ Code Quality

## lint: Run linter
lint:
	@echo "$(CYAN)Running linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)Linting complete$(NC)"

## fmt: Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .
	@echo "$(GREEN)Formatting complete$(NC)"

## vet: Run go vet
vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Vet complete$(NC)"

##@ Database

## migrate-up: Run database migrations
migrate-up:
	@echo "$(CYAN)Running migrations...$(NC)"
	@psql $(DB_DSN) -f migrations/001_create_todos_table.sql
	@echo "$(GREEN)Migrations complete$(NC)"

## migrate-down: Rollback database migrations
migrate-down:
	@echo "$(CYAN)Rolling back migrations...$(NC)"
	@psql $(DB_DSN) -c "DROP TABLE IF EXISTS todos CASCADE;"
	@echo "$(GREEN)Rollback complete$(NC)"

## db-create: Create database
db-create:
	@echo "$(CYAN)Creating database...$(NC)"
	@psql "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD)" -c "CREATE DATABASE $(DB_NAME);"
	@echo "$(GREEN)Database created$(NC)"

## db-drop: Drop database
db-drop:
	@echo "$(CYAN)Dropping database...$(NC)"
	@psql "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD)" -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "$(GREEN)Database dropped$(NC)"

##@ Docker

## docker-up: Start Docker services
docker-up:
	@echo "$(CYAN)Starting Docker services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Docker services started$(NC)"

## docker-down: Stop Docker services
docker-down:
	@echo "$(CYAN)Stopping Docker services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Docker services stopped$(NC)"

##@ Setup

## install-tools: Install required tools
install-tools:
	@echo "$(CYAN)Installing tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)Tools installed$(NC)"

## install-hooks: Install git hooks
install-hooks:
	@echo "$(CYAN)Installing git hooks...$(NC)"
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Git hooks installed$(NC)"

## deps: Download dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

##@ CI

## ci: Run all CI checks (format, lint, test)
ci: fmt vet lint test
	@echo "$(GREEN)All CI checks passed!$(NC)"
