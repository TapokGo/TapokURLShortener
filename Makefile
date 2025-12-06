.DEFAULT_GOAL := help
PROJECT_NAME := url_shortener
BIN_DIR := bin

.PHONY: test lint build run clean help

lint:
	@echo "--------------------------------"
	@echo "Starting linters..."
	@golangci-lint run

test: 
	@echo "--------------------------------"
	@echo "Starting tests..."
	@go test -v ./...
	@echo "Tests complete"

build: lint test
	@echo "--------------------------------"
	@echo "Building $(PROJECT_NAME)..."
	@go build -o ./$(BIN_DIR)/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)
	@echo "Build complete: $(BIN_DIR)/$(PROJECT_NAME)"

run: build	
	@echo "--------------------------------"
	@echo "Starting project..."
	@./$(BIN_DIR)/$(PROJECT_NAME)

clean:
	@echo "--------------------------------"
	@echo "Cleaning up..."
	@rm -rf ./$(BIN_DIR)
	@rm -rf ./storage
	@echo "Cleanup complete."
	
help:
	@echo "Available commands:"
	@echo "		make lint		- runs linters check"
	@echo "		make test		- runs all go test"
	@echo "		make build		- compiles the Go application(also start lint and tests)"
	@echo "		make run		- build and run application"
	@echo "		make clean		- removes complied binaries and other generated files"
	@echo "		make help		- display this help message"


