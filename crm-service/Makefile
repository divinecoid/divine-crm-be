.PHONY: help build run test clean migrate seed docker-build docker-up docker-down

# Variables
APP_NAME=divine-crm
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_PATH)

run: ## Run the application
	$(GO) run $(MAIN_PATH)

dev: ## Run with live reload (requires air)
	air

test: ## Run tests
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	$(GO) tool cover -html=coverage.out

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete"

migrate: ## Run database migrations
	@echo "Running migrations..."
	@$(GO) run scripts/migrate.go up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@$(GO) run scripts/migrate.go down

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@$(GO) run scripts/seed.go

lint: ## Run linter
	golangci-lint run

format: ## Format code
	$(GO) fmt ./...
	goimports -w .

docker-build: ## Build Docker image
	$(DOCKER) build -t $(APP_NAME):latest .

docker-up: ## Start Docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-clean: ## Clean Docker resources
	$(DOCKER_COMPOSE) down -v
	$(DOCKER) system prune -f

db-create: ## Create database
	createdb divine_crm

db-drop: ## Drop database
	dropdb divine_crm

db-reset: db-drop db-create migrate seed ## Reset database

generate-swagger: ## Generate Swagger documentation
	swag init -g $(MAIN_PATH)

install-tools: ## Install development tools
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

setup: deps install-tools ## Setup development environment
	@echo "Development environment ready!"

prod-build: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -ldflags '-s -w' -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Default target
.DEFAULT_GOAL := help
