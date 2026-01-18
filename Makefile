MAKEFILE_VERSION := 1.0.0

.PHONY: help build run run-dev test lint docker-build docker-run clean dev prod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	@echo "Building service-template..."
	@go build -o service-template .

run: ## Run the application locally (production mode)
	@echo "Running service-template (production)..."
	@DISPLAY_MESSAGE="🚀 Production Service Running!" \
		ENVIRONMENT=production \
		VERSION=1.0.0 \
		PORT=8080 \
		go run main.go

run-dev: ## Run the application locally (development mode)
	@echo "Running service-template (development)..."
	@DISPLAY_MESSAGE="🛠️ Development/Staging Environment - Testing in Progress" \
		ENVIRONMENT=development \
		VERSION=dev-latest \
		PORT=8181 \
		go run main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t service-template:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 \
		-e DISPLAY_MESSAGE="🚀 Production Service Running!" \
		-e ENVIRONMENT=production \
		service-template:latest

docker-run-dev: ## Run Docker container in dev mode
	@echo "Running Docker container (dev)..."
	@docker run -p 8081:8080 \
		-e DISPLAY_MESSAGE="🛠️ Development Environment" \
		-e ENVIRONMENT=development \
		service-template:latest

prod: ## Start production environment with docker-compose
	@echo "Starting production environment..."
	@docker-compose up -d
	@echo "Production service available at http://localhost:8080"

dev: ## Start development environment with docker-compose
	@echo "Starting development environment..."
	@docker-compose -f docker-compose.dev.yml up -d
	@echo "Development service available at http://localhost:8081"

stop: ## Stop all docker-compose services
	@echo "Stopping services..."
	@docker-compose down
	@docker-compose -f docker-compose.dev.yml down

logs: ## Show docker-compose logs
	@docker-compose logs -f

logs-dev: ## Show docker-compose dev logs
	@docker-compose -f docker-compose.dev.yml logs -f

clean: ## Clean build artifacts and containers
	@echo "Cleaning up..."
	@rm -f service-template coverage.out
	@docker-compose down -v
	@docker-compose -f docker-compose.dev.yml down -v

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

.DEFAULT_GOAL := help
