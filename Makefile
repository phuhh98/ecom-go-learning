# Makefile for E-commerce Go project with cross-platform support

# Variables
BIN_DIR = bin
CMD_DIR = cmd
DOCKER_COMPOSE = docker/docker-compose.yml
DOCKER = podman

# Detect OS
ifeq ($(OS),Windows_NT)
    # Windows
    BINARY_EXT = .exe
    RM = rmdir /s /q
    MKDIR = mkdir
else
    # macOS/Linux
    BINARY_EXT =
    RM = rm -rf
    MKDIR = mkdir -p
endif

# Build binary executables
build: swagger-gen
	@$(MKDIR) $(BIN_DIR) || true
	go build -o $(BIN_DIR)/api$(BINARY_EXT) $(CMD_DIR)/api/main.go
	go build -o $(BIN_DIR)/worker$(BINARY_EXT) $(CMD_DIR)/worker/main.go

# Development with hot reload
# 1. Start dependencies (PostgreSQL, Redis, RabbitMQ) in Docker
deps-up:
	$(DOCKER) compose -f $(DOCKER_COMPOSE) up -d

# 2. Run API service with OS detection
run-api:
ifeq ($(OS),Windows_NT)
	air -c $(CMD_DIR)/api/.air.toml
else
	air -c $(CMD_DIR)/api/.air.toml --build.cmd "go build -o ./tmp/api ./cmd/api" --build.bin="./tmp/api"
endif

# Run API in production mode
run-api-prod: build
	@echo "Starting API in production mode..."
	ENV=production $(BIN_DIR)/api$(BINARY_EXT)

# 3. Run worker service with OS detection
run-worker:
ifeq ($(OS),Windows_NT)
	air -c $(CMD_DIR)/worker/.air.toml
else
	air -c $(CMD_DIR)/worker/.air.toml --build.cmd "go build -o ./tmp/worker ./cmd/worker" --build.bin="./tmp/worker"
endif

# Stop dependencies
deps-down:
	$(DOCKER) compose -f $(DOCKER_COMPOSE) down

# Format code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	$(RM) $(BIN_DIR)

# install swag cli
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# auto generate swagger docs
swagger-gen:
	swag init --parseDependency --generalInfo cmd/api/main.go

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the project binaries"
	@echo "  deps-up        - Start all dependencies"
	@echo "  run-api        - Run the API service"
	@echo "  run-api-prod   - Run the API service (production)"
	@echo "  run-worker     - Run the worker service"
	@echo "  deps-down      - Stop all dependencies"
	@echo "  fmt            - Format code"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  install-swag   - Install swag CLI"
	@echo "  swagger-gen    - Generate Swagger documentation"

.PHONY: build deps-up run-api run-api-prod run-worker deps-down fmt deps clean help install-swag swagger-gen