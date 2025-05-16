.PHONY: build run test clean generate docker-build docker-run install-tools help

# Build variables
BINARY_NAME=gateway-service
MAIN_FILE=./cmd/app/main.go

# Docker variables
DOCKER_IMAGE=gateway-service
DOCKER_TAG=latest

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOMOD=$(GO) mod

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Generate code
generate:
	$(GO) generate ./...

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	docker-compose up

# Install development tools
install-tools:
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/go-delve/delve/cmd/dlv

# Help command
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build files"
	@echo "  generate      - Generate code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message" 