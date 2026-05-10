# Makefile for Toucan LMS

# Variables
BINARY_NAME=toucan
DOCKER_IMAGE=toucan
TAG?=latest
GO_CMD=go
PNPM_CMD=pnpm
UI_DIR=ui
BIN_DIR=bin

# Default target
.PHONY: all
all: build

# Help target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make install         - Install both backend and frontend dependencies"
	@echo "  make build           - Build both backend and frontend"
	@echo "  make dev             - Run both backend and frontend in development mode"
	@echo "  make test            - Run both backend and frontend tests"
	@echo "  make lint            - Lint both backend and frontend"
	@echo "  make tidy            - Run go mod tidy"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make run-backend     - Run only the backend"
	@echo "  make run-frontend    - Run only the frontend"
	@echo "  make docker-build    - Build production Docker image (including frontend)"
	@echo "  make docker-headless - Build headless Docker image (API only)"

# Install dependencies
.PHONY: install
install: install-backend install-frontend

.PHONY: install-backend
install-backend:
	$(GO_CMD) mod download

.PHONY: install-frontend
install-frontend:
	cd $(UI_DIR) && $(PNPM_CMD) install

# Build targets
.PHONY: build
build: build-backend build-frontend

.PHONY: build-backend
build-backend:
	@echo "Building backend..."
	@mkdir -p $(BIN_DIR)
	$(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) cmd/toucan/main.go

.PHONY: build-frontend
build-frontend:
	@echo "Building frontend..."
	cd $(UI_DIR) && $(PNPM_CMD) install && $(PNPM_CMD) run build --outDir ../public

# Run targets
.PHONY: run-backend
run-backend:
	@echo "Running backend..."
	$(GO_CMD) run cmd/toucan/main.go

.PHONY: run-frontend
run-frontend:
	@echo "Running frontend..."
	cd $(UI_DIR) && $(PNPM_CMD) run dev

# Run both backend and frontend in parallel
.PHONY: dev
dev:
	@echo "Starting development environment..."
	@make -j 2 run-backend run-frontend

# Test targets
.PHONY: test
test: test-backend test-frontend

.PHONY: test-backend
test-backend:
	@echo "Running backend tests..."
	$(GO_CMD) test ./...

.PHONY: test-frontend
test-frontend:
	@echo "Running frontend tests..."
	cd $(UI_DIR) && $(PNPM_CMD) run test:e2e

# Maintenance
.PHONY: tidy
tidy:
	$(GO_CMD) mod tidy

.PHONY: lint
lint: lint-backend lint-frontend

.PHONY: lint-backend
lint-backend:
	@echo "Linting backend..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping backend lint"; \
	fi

.PHONY: lint-frontend
lint-frontend:
	@echo "Linting frontend..."
	cd $(UI_DIR) && $(PNPM_CMD) run lint

# Clean
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	rm -rf $(UI_DIR)/dist
	rm -rf public/

# Docker build targets
.PHONY: docker-build
docker-build:
	@echo "Building production Docker image (tag: $(TAG))..."
	docker build --build-arg HEADLESS=false -t $(DOCKER_IMAGE):$(TAG) .

.PHONY: docker-headless
docker-headless:
	@echo "Building headless Docker image (tag: $(TAG)-headless)..."
	docker build --build-arg HEADLESS=true -t $(DOCKER_IMAGE):$(TAG)-headless .
